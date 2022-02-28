package repository

import (
	"context"
	"github.com/jackc/pgx"
	"go-musthave-diploma-tpl/internal/config"
)

func NewRepository(cfg *config.ServerConfiguration) (Storage, error) {
	//TODO вернуть ли ошибку???
	return NewDBConnect(cfg.DatabaseDSN)
}

type DBRepo struct {
	conn *pgx.Conn
}

func NewDBConnect(databaseDSN string) (*DBRepo, error) {
	conn, err := pgx.Connect(context.Background(), databaseDSN)
	if err != nil {
		return nil, err
	}

	pgRepo := DBRepo{
		conn: conn,
	}

	m, err := RunMigration(databaseDSN)
	if err != nil && !m {
		return nil, err
	}
	//TODO стоит ли добавить что-нибудь в таблицы???

	return &pgRepo, nil
}
