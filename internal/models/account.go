package models

import (
	"github.com/ShiraazMoollatjie/goluhn"
)

// Account счет пользователя с аккумулированным балансом и суммой списаний
type Account struct {
	AccountID   string
	UID         string
	Balance     float32
	Withdrawals float32
}

func NewAccount(userUID string, balance float32, withdrawals float32) Account {
	return Account{
		AccountID:   goluhn.Generate(16),
		UID:         userUID,
		Balance:     balance,
		Withdrawals: withdrawals,
	}
}
