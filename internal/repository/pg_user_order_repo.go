package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type UserOrderRepository struct {
	conn *pgx.Conn
}

func (u UserOrderRepository) AddUserOrder(ctx context.Context, userUID string, order string) {

}

func NewUserOrderRepository(conn *pgx.Conn) *UserOrderRepository {
	return &UserOrderRepository{conn: conn}
}
