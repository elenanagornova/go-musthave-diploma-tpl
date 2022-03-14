package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"go-musthave-diploma-tpl/internal/gophermart"
	"go-musthave-diploma-tpl/internal/models"
)

var queryGetUser = "select user_uid, login, password from gophermart.users where login=$1"
var queryAddUser = "insert into gophermart.users(user_uid, login, password) values($1, $2, $3)"

type UserRepository struct {
	conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{conn: conn}
}

func (u UserRepository) GetUser(ctx context.Context, login string) (models.User, error) {
	var user models.User
	result := u.conn.QueryRow(ctx, queryGetUser, login)
	if err := result.Scan(&user.UID, &user.Login, &user.Password); err != nil {
		if err.Error() == "no rows in result set" {
			return models.User{}, gophermart.ErrUserNotFound
		}
		return models.User{}, err
	}
	return user, nil
}

func (u UserRepository) AddUser(ctx context.Context, user models.User) error {
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, error := u.conn.Exec(ctx, queryAddUser, user.UID, user.Login, user.Password)
	if error != nil {
		return error
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
