package controller

import (
	"context"
	"encoding/json"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go-musthave-diploma-tpl/gen/gophermart"
	"go-musthave-diploma-tpl/internal/gophermart"
	"go-musthave-diploma-tpl/internal/pkg/auth"
	"net/http"
	"strconv"
)

func NewRouter(context context.Context, service *gophermart.Gophermart) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(JwtMiddleware)
	//
	r.Post("/api/user/register", UserRegister(context, service))
	r.Post("/api/user/login", LoginUser(context, service))
	r.Post("/api/user/orders", UploadUserOrder(context, service))
	r.Get("/api/user/orders", GetUserOrders(context, service))
	r.Get("/api/user/balance", GetUserBalance(context, service))
	r.Post("/api/user/balance/withdraw", WithdrawBalance(context, service))
	r.Get("/ping", CheckConn(context, service))
	return r
}

func CheckConn(ctx context.Context, service *gophermart.Gophermart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func WithdrawBalance(ctx context.Context, service *gophermart.Gophermart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func GetUserBalance(ctx context.Context, service *gophermart.Gophermart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func GetUserOrders(ctx context.Context, service *gophermart.Gophermart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func UploadUserOrder(ctx context.Context, service *gophermart.Gophermart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var orderNum int
		err := json.NewDecoder(r.Body).Decode(&orderNum)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		strOrderNum := strconv.Itoa(orderNum)
		if goluhn.Validate(strOrderNum) != nil {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		login := userLoginFromRequest(r)
		addOrderErr := service.AddOrder(ctx, login, strOrderNum)
		if addOrderErr != nil {
			if addOrderErr == gophermart.ErrOrderOwnedByAnotherUser {
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			} else if addOrderErr == gophermart.ErrOrderExists {
				w.WriteHeader(http.StatusOK)
				return
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		token, _ := auth.CreateToken(login)
		http.SetCookie(w, &token)
		w.WriteHeader(http.StatusAccepted)
	}
}
func userLoginFromRequest(r *http.Request) string {
	uid := r.Context().Value(UserCtxKey)
	if uid == nil {
		return ""
	}
	if userID, ok := uid.(string); ok {
		return userID
	}
	return ""
}

//LoginUser аутентификация юзера
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
			if gophermart.IsPgUniqueViolationError(error) {
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		//	var pgErr *pgconn.PgError
		//	if errors.As(error, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		//		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		//		return
		//	}
		//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		//	return
		//}
		//TODO: с миддлваре сделать?
		token, _ := auth.CreateToken(registerRequest.Login)
		http.SetCookie(w, &token)
		w.WriteHeader(http.StatusOK)
	}
}
