package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"go-musthave-diploma-tpl/internal/models"
)

var queryGetOrders = "select order_num, sum, processed_at from withdrawals"
var queryAddWithdrawal = "insert into gophermart.withdrawals(order_num, sum, processed_at) values($1, $2, $3)"

type WithdrawalRepository struct {
	conn *pgx.Conn
}

func (w WithdrawalRepository) AddWithdraw(ctx context.Context, withdrawal models.Withdrawal) error {
	tx, err := w.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, withdrawErr := w.conn.Exec(ctx, queryAddWithdrawal, withdrawal.OrderNum, withdrawal.Sum, withdrawal.ProcessedAt)
	if withdrawErr != nil {
		return withdrawErr
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (w WithdrawalRepository) GetWithdrawals(ctx context.Context) []models.Withdrawal {
	var result []models.Withdrawal
	rows, err := w.conn.Query(ctx, queryGetOrders)
	if err != nil {
		return nil
	}
	for rows.Next() {
		var withdrawal models.Withdrawal
		err = rows.Scan(&withdrawal.OrderNum, &withdrawal.Sum, &withdrawal.ProcessedAt)
		if err != nil {
			return nil
		}
		result = append(result, withdrawal)
	}
	err = rows.Err()
	if err != nil {
		return nil
	}
	return result
}

func NewWithdrawalRepository(conn *pgx.Conn) *WithdrawalRepository {
	return &WithdrawalRepository{conn: conn}
}
