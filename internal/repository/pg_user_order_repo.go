package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"go-musthave-diploma-tpl/internal/gophermart"
	"go-musthave-diploma-tpl/internal/models"
)

var queryAddOrder = "insert into gophermart.orders(user_uid, order_num, uploaded_at, status, accrual) values($1, $2, $3, $4, $5)"
var queryGetOrder = "select user_uid, order_num, uploaded_at, status, accrual  from gophermart.orders where order_num=$1"
var queryGetOrdersByUserUID = "select user_uid, order_num, uploaded_at, status, accrual  from gophermart.orders where user_uid=$1 order by uploaded_at desc"
var queryGetAllOrders = "select user_uid, order_num, uploaded_at, status, accrual  from gophermart.orders order by uploaded_at desc"
var queryGetNewAndProcessingOrders = "select user_uid, order_num, uploaded_at, status, accrual  from gophermart.orders where status='NEW' OR status='new' OR status='PROCESSING'"
var queryUpdateOrderState = "update gophermart.orders set status=$1, accrual=$2, retry_count=$3 where order_num=$4 and status!=$1"

type UserOrderRepository struct {
	conn *pgx.Conn
}

func (u UserOrderRepository) GetAllUserOrders(ctx context.Context) []models.Order {
	var result []models.Order
	rows, err := u.conn.Query(ctx, queryGetAllOrders)
	if err != nil {
		return nil
	}
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.UID, &order.OrderID, &order.UploadedAt, &order.Status, &order.Accrual)
		if err != nil {
			return nil
		}
		result = append(result, order)
	}
	err = rows.Err()
	if err != nil {
		return nil
	}
	return result
}

func (u UserOrderRepository) UpdateOrdersStateFromAccrual(ctx context.Context, orders []gophermart.Result) error {
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	for _, value := range orders {
		_, err = tx.Exec(ctx, queryUpdateOrderState, value.Status, value.Accrual, value.RetryCount, value.OrderID)
		if err != nil {
			return err
		}
		tx.Commit(ctx)
	}
	return nil
}

func (u UserOrderRepository) GetNewAndProcessingOrders(ctx context.Context) []models.Order {
	var result []models.Order
	rows, err := u.conn.Query(ctx, queryGetNewAndProcessingOrders)
	if err != nil {
		return nil
	}
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.UID, &order.OrderID, &order.UploadedAt, &order.Status, &order.Accrual)
		if err != nil {
			return nil
		}
		result = append(result, order)
	}
	err = rows.Err()
	if err != nil {
		return nil
	}
	return result
}

func (u UserOrderRepository) GetUserOrders(ctx context.Context, userUID string) []models.Order {
	var result []models.Order
	rows, err := u.conn.Query(ctx, queryGetOrdersByUserUID, userUID)
	if err != nil {
		return nil
	}
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.UID, &order.OrderID, &order.UploadedAt, &order.Status, &order.Accrual)
		if err != nil {
			return nil
		}
		result = append(result, order)
	}
	err = rows.Err()
	if err != nil {
		return nil
	}
	return result
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
