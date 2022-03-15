package gophermart

import (
	"context"
	"go-musthave-diploma-tpl/internal/models"
)

type UserRepo interface {
	GetUser(ctx context.Context, login string) (models.User, error)
	AddUser(ctx context.Context, entity models.User) error
}

type UserAccountRepo interface {
	AddAccount(ctx context.Context, account models.Account) error
	GetAccount(ctx context.Context, userID string) (models.Account, error)
	RefillAmount(ctx context.Context, userID string, diff float32) error
	WithdrawalAmount(ctx context.Context, userID string, diff float32, orderNum string) error
}

type UserOrderRepo interface {
	AddUserOrder(ctx context.Context, order models.Order) error
	GetOrder(ctx context.Context, orderNum string) (models.Order, error)
	GetUserOrders(ctx context.Context, userUID string) []models.Order
	GetNewAndProcessingOrders(ctx context.Context) []models.Order
	UpdateOrdersStateFromAccrual(ctx context.Context, orders []Result) error
	GetAllUserOrders(ctx context.Context) []models.Order
}

type WithdrawalRepo interface {
	GetWithdrawals(ctx context.Context) []models.Withdrawal
}
