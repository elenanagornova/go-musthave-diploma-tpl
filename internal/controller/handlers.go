package controller

import (
	"go-musthave-diploma-tpl/internal/gophermart"
	"net/http"
)

func RegisterUser(service *gophermart.Loyalty) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func LoginUser(service *gophermart.Loyalty) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func UploadUserOrder(service *gophermart.Loyalty) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func GetUserOrders(service *gophermart.Loyalty) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func GetUserBalance(service *gophermart.Loyalty) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func WithdrawBalance(service *gophermart.Loyalty) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
