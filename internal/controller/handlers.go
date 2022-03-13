package controller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go-musthave-diploma-tpl/gen/gophermart"
	"go-musthave-diploma-tpl/internal/gophermart"
	"go-musthave-diploma-tpl/internal/pkg/auth"
	"net/http"
)

func NewRouter(context context.Context, service *gophermart.Gophermart) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	//
	r.Post("/api/user/register", UserRegister(context, service))
	r.Post("/api/user/login", LoginUser(context, service))
	//r.Post("/api/user/orders", controller.UploadUserOrder(service))
	//r.Get("/api/user/orders", controller.GetUserOrders(service))
	//r.Get("/api/user/balance", controller.GetUserBalance(service))
	//r.Post("/api/user/balance/withdraw", controller.WithdrawBalance(service))
	//r.Get("/ping", controller.CheckConn(service))
	return r
}

func LoginUser(context context.Context, service *gophermart.Gophermart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userLoginRequest Openapi.UserLoginRequest
		err := json.NewDecoder(r.Body).Decode(&userLoginRequest)
		if err != nil || userLoginRequest.Login == "" || userLoginRequest.Password == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := service.LoginUser(context, userLoginRequest.Login, userLoginRequest.Password); err != nil {
			if err == gophermart.ErrAuth {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if err == gophermart.ErrUserNotFound {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		}
		token, _ := auth.CreateToken(userLoginRequest.Login)
		http.SetCookie(w, &token)
		w.WriteHeader(http.StatusOK)
	}
}
func UserRegister(context context.Context, service *gophermart.Gophermart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var registerRequest Openapi.UserRegisterRequest
		err := json.NewDecoder(r.Body).Decode(&registerRequest)
		if err != nil || registerRequest.Login == "" || registerRequest.Password == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if error := service.RegisterUser(context, registerRequest.Login, registerRequest.Password); error != nil {
			var pgErr *pgconn.PgError
			if errors.As(error, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		//TODO: с миддлваре сделать?
		token, _ := auth.CreateToken(registerRequest.Login)
		http.SetCookie(w, &token)
		w.WriteHeader(http.StatusOK)
	}
}
