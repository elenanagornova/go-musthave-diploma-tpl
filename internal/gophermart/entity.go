package gophermart

import "errors"

var (
	ErrUserExists              = errors.New("user already exists")
	ErrUserNotFound            = errors.New("user not found")
	ErrAuth                    = errors.New("invalid login or password")
	ErrOrderExists             = errors.New("order already exists")
	ErrOrderOwnedByAnotherUser = errors.New("order uploaded by another user")
	ErrInvalidOrderFormat      = errors.New("order format error")
	ErrNoMoney                 = errors.New("no money")
	ErrOrderNotFound           = errors.New("order not found for user")
)
