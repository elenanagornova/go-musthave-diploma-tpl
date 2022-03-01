package repository

type Storage interface {
	Ping() error
}
