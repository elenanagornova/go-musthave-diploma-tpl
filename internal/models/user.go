package models

import (
	"github.com/google/uuid"
)

type User struct {
	UID      string
	Login    string
	Password string //хеш пароля
}

func NewUser(login string, password string) (User, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return User{}, err
	}
	user := User{
		UID:      uid.String(),
		Login:    login,
		Password: password,
	}
	return user, nil
}
