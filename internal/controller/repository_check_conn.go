package controller

import (
	"go-musthave-diploma-tpl/internal/gophermart"
	"net/http"
)

func CheckConn(service *gophermart.Loyalty) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := service.Repo.Ping(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
