package gophermart

import (
	"context"
	"go-musthave-diploma-tpl/internal/models"
)

// UserRepo для работы с сущностью User
type UserRepo interface {
	GetUser(ctx context.Context, login string) (models.User, error)
	AddUser(ctx context.Context, entity models.User) error
}

type UserAccountRepo interface {
	AddAccount(ctx context.Context, account models.Account) error
	GetAccount(ctx context.Context, userID string) (models.Account, error)
	RefillAmount(ctx context.Context, userID string, diff float32) error
	WithdrawalAmount(ctx context.Context, userID string, diff float32) error
}

type UserOrderRepo interface {
	AddUserOrder(ctx context.Context, userUID string, order string)
}
