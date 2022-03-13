package gophermart

import (
	"context"
	"go-musthave-diploma-tpl/internal/models"
	"go-musthave-diploma-tpl/internal/pkg/hasher"
)

type Gophermart struct {
	Addr            string
	UserRepo        UserRepo
	UserAccountRepo UserAccountRepo
}

func (g Gophermart) RegisterUser(context context.Context, login string, password string) error {
	hashedPassword, err := hasher.HashAndSalt([]byte(password))
	if err != nil {
		return err
	}

	newUser, err := models.NewUser(login, hashedPassword)
	if err != nil {
		return err
	}

	if err := g.UserRepo.AddUser(context, newUser); err != nil {
		return err
	}

	newAccount := models.NewAccount(newUser.UID, 0, 0)
	if err := g.UserAccountRepo.AddAccount(context, newAccount); err != nil {
		return err
	}
	return nil
}

func NewGophermart(addr string, userRepo UserRepo, userAccountRepo UserAccountRepo) *Gophermart {
	return &Gophermart{Addr: addr, UserRepo: userRepo, UserAccountRepo: userAccountRepo}
}
