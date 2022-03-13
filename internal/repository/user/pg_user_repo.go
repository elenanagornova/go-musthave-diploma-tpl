package user

import (
	"context"
	"github.com/jackc/pgx/v4"
	"go-musthave-diploma-tpl/internal/models"
)

var queryGetUser = "select uid, login, password from gophermart.users where login=$1"
var queryAddUser = "insert into gophermart.users(user_uid, login, password) values($1, $2, $3)"

type UserRepository struct {
	conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{conn: conn}
}

func (u UserRepository) GetUser(ctx context.Context, login string) (models.User, error) {
	return models.User{}, nil
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
