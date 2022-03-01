package controller

import (
	"go-musthave-diploma-tpl/internal/loyalty_system"
	"net/http"
)

func RegisterUser(service *loyalty_system.Loyalty) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func LoginUser(service *loyalty_system.Loyalty) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func UploadUserOrder(service *loyalty_system.Loyalty) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func GetUserOrders(service *loyalty_system.Loyalty) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func GetUserBalance(service *loyalty_system.Loyalty) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func WithdrawBalance(service *loyalty_system.Loyalty) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
