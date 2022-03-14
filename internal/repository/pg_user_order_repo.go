package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"go-musthave-diploma-tpl/internal/gophermart"
	"go-musthave-diploma-tpl/internal/models"
)

var queryAddOrder = "insert into gophermart.orders(user_uid, order_num, uploaded_at, status, accrual) values($1, $2, $3, $4, $5)"
var queryGetOrder = "select user_uid, order_num, uploaded_at, status, accrual  from gophermart.orders where order_num=$1"

type UserOrderRepository struct {
	conn *pgx.Conn
}

func (u UserOrderRepository) GetOrder(ctx context.Context, orderNum string) (models.Order, error) {
	var order models.Order
	result := u.conn.QueryRow(ctx, queryGetOrder, orderNum)
	if err := result.Scan(&order.UID, &order.OrderID, &order.UploadedAt, &order.Status, &order.Accrual); err != nil {
		if err.Error() == "no rows in result set" {
			return models.Order{}, gophermart.ErrOrderExists
		}
		return models.Order{}, err
	}
	return order, nil
}

func (u UserOrderRepository) AddUserOrder(ctx context.Context, order models.Order) error {
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, error := u.conn.Exec(ctx, queryAddOrder, order.UID, order.OrderID, order.UploadedAt, models.NEW, 0)
	if error != nil {
		return error
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func NewUserOrderRepository(conn *pgx.Conn) *UserOrderRepository {
	return &UserOrderRepository{conn: conn}
}
