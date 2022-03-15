package gophermart

import (
	"context"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	Accrual "go-musthave-diploma-tpl/gen/accrual"
	Openapi "go-musthave-diploma-tpl/gen/gophermart"
	"go-musthave-diploma-tpl/internal/models"
	"go-musthave-diploma-tpl/internal/pkg/hasher"
	"go-musthave-diploma-tpl/internal/providers"
	"time"
)

const accrualDefaultRetryInterval = 50 * time.Millisecond

type Gophermart struct {
	Addr                 string
	Accrual              Accrual.Client
	UserRepo             UserRepo
	UserAccountRepo      UserAccountRepo
	UserOrderRepo        UserOrderRepo
	WithdrawalRepo       WithdrawalRepo
	accrualProvider      providers.Provider
	accrualRetryInterval time.Duration
}

func (g Gophermart) RegisterUser(context context.Context, login string, password string) error {
	hashedPassword, err := hasher.HashAndSalt([]byte(password))
	if err != nil {
		return err
	}

	newUser, err := models.NewUser(login, hashedPassword)
	if err != nil {
		return err
	}

	if err := g.UserRepo.AddUser(context, newUser); err != nil {
		return err
	}

	newAccount := models.NewAccount(newUser.UID, 0, 0)
	if err := g.UserAccountRepo.AddAccount(context, newAccount); err != nil {
		return err
	}
	return nil
}

func (g Gophermart) LoginUser(ctx context.Context, login string, password string) error {
	user, err := g.UserRepo.GetUser(ctx, login)
	if err != nil {
		if err == ErrUserNotFound {
			return ErrUserNotFound
		}
		return err
	}

	if !hasher.CheckPassword(user.Password, []byte(password)) {
		return ErrAuth
	}
	return nil
}

func (g Gophermart) AddOrder(ctx context.Context, login string, orderNum string) error {
	user, err := g.UserRepo.GetUser(ctx, login)
	if err != nil {
		return err
	}
	order := models.NewOrder(user.UID, orderNum)
	if errDB := g.UserOrderRepo.AddUserOrder(ctx, order); errDB != nil {
		if IsPgUniqueViolationError(errDB) {
			order, err := g.UserOrderRepo.GetOrder(ctx, orderNum)
			if err != nil {
				return err
			}
			if order.UID == user.UID {
				return ErrOrderExists
			} else {
				return ErrOrderOwnedByAnotherUser
			}
		}
	}
	return nil
}

func (g Gophermart) GetUserOrders(ctx context.Context, login string) []models.Order {
	user, err := g.UserRepo.GetUser(ctx, login)
	if err != nil {
		return nil
	}
	return g.UserOrderRepo.GetUserOrders(ctx, user.UID)
}

func (g Gophermart) GetUserBalance(ctx context.Context, login string) (Openapi.GetUserBalanceResponse, error) {
	user, err := g.UserRepo.GetUser(ctx, login)
	if err != nil {
		return Openapi.GetUserBalanceResponse{}, err
	}
	account, err := g.UserAccountRepo.GetAccount(ctx, user.UID)
	if err != nil {
		return Openapi.GetUserBalanceResponse{}, err
	}
	return Openapi.GetUserBalanceResponse{Current: account.Balance, Withdrawn: account.Withdrawals}, nil
}

func (g Gophermart) WithDrawUserBalance(ctx context.Context, login string, withdrawRequest Openapi.UserBalanceWithdrawRequest) error {
	return g.UserAccountRepo.WithdrawalAmount(ctx, login, withdrawRequest.Sum, withdrawRequest.Order)
}

func (g Gophermart) GetWithdrawals(ctx context.Context) []models.Withdrawal {
	return g.WithdrawalRepo.GetWithdrawals(ctx)

}
func NewGophermart(addr string, userRepo UserRepo, userAccountRepo UserAccountRepo, userOrderRepo UserOrderRepo, withdrawalRepo WithdrawalRepo, accrualProvider providers.Provider) *Gophermart {
	return &Gophermart{Addr: addr, UserRepo: userRepo, UserAccountRepo: userAccountRepo, UserOrderRepo: userOrderRepo, WithdrawalRepo: withdrawalRepo, accrualProvider: accrualProvider, accrualRetryInterval: accrualDefaultRetryInterval}
}

func IsPgUniqueViolationError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
}
