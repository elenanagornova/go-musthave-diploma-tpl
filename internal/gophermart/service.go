package gophermart

import "go-musthave-diploma-tpl/internal/repository"

// Loyalty service
type Loyalty struct {
	Addr string
	Repo repository.Storage
}

// New Loyalty instance
func New(addr string, storage repository.Storage) *Loyalty {
	return &Loyalty{
		Addr: addr,
		Repo: storage,
	}
}
