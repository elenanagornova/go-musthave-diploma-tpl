package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"go-musthave-diploma-tpl/internal/models"
)

type UserAccountRepository struct {
	conn *pgx.Conn
}

func NewUserAccountRepository(conn *pgx.Conn) *UserAccountRepository {
	return &UserAccountRepository{conn: conn}
}

var queryAddAccount = "insert into gophermart.account_info(user_uid, balance, withdrawal) values($1, $2, $3)"

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
	panic("implement me")
}

func (u UserAccountRepository) RefillAmount(ctx context.Context, userID string, diff float32) error {
	panic("implement me")
}

func (u UserAccountRepository) WithdrawalAmount(ctx context.Context, userID string, diff float32) error {
	panic("implement me")
}

func (u UserAccountRepository) Close() error {
	panic("implement me")
}
