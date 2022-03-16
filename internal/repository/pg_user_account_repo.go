package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"go-musthave-diploma-tpl/internal/gophermart"
	"go-musthave-diploma-tpl/internal/models"
	"strconv"
)

type UserAccountRepository struct {
	conn *pgx.Conn
}

func NewUserAccountRepository(conn *pgx.Conn) *UserAccountRepository {
	return &UserAccountRepository{conn: conn}
}

var queryAddAccount = "insert into gophermart.account_info(user_uid, balance, withdrawal) values($1, $2, $3)"
var queryGetBalance = "select balance, withdrawal from gophermart.account_info where user_uid=$1"
var queryGetUserBalanceOrder = "SELECT u.user_uid, a.balance, a.withdrawal, o.order_num FROM gophermart.users u INNER JOIN  gophermart.account_info a on a.user_uid = u.user_uid JOIN gophermart.orders o on a.user_uid = o.user_uid WHERE u.login=$1"
var queryUpdateByWithdrawal = "update gophermart.account_info set balance=$1, withdrawal=$2 where user_uid=$3"
var queryUpdateBalance = "update gophermart.account_info set balance=($1+balance) where user_uid=(select user_uid from gophermart.orders where order_num=$2)"

type UserBalanceWithOrder struct {
	userUID    string
	balance    float32
	withdrawal float32
	orderNum   string
}

func (u UserAccountRepository) AddAccount(ctx context.Context, account models.Account) error {
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, error := u.conn.Exec(ctx, queryAddAccount, account.UID, account.Balance, account.Withdrawals)
	if error != nil {
		return error
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (u UserAccountRepository) GetAccount(ctx context.Context, userID string) (models.Account, error) {
	var account models.Account
	result := u.conn.QueryRow(ctx, queryGetBalance, userID)
	if err := result.Scan(&account.Balance, &account.Withdrawals); err != nil {
		if err.Error() == "no rows in result set" {
			return models.Account{}, gophermart.ErrUserNotFound
		}
		return models.Account{}, err
	}
	return account, nil
}
func (u UserAccountRepository) WithdrawalAmount(ctx context.Context, login string, diff float32, orderNum string) error {
	var user = UserBalanceWithOrder{}
	result := u.conn.QueryRow(ctx, queryGetUserBalanceOrder, login)
	if err := result.Scan(&user.userUID, &user.balance, &user.withdrawal, &user.orderNum); err != nil {
		return err
	}
	if user.balance < diff {
		return gophermart.ErrNoMoney
	}
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, error := u.conn.Exec(ctx, queryUpdateByWithdrawal, user.balance-diff, user.withdrawal, user.userUID)
	if error != nil {
		return error
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
func (u UserAccountRepository) RefillAmount(ctx context.Context, userID string, diff float32) error {
	panic("implement me")
}
func (u UserAccountRepository) UpdateBalance(ctx context.Context, order gophermart.Result) error {
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(ctx, queryUpdateBalance, order.Accrual, strconv.Itoa(order.OrderID))
	if err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
func (u UserAccountRepository) Close() error {
	panic("implement me")
}
