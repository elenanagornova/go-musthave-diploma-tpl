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
	r.Get("/api/user/withdrawals", GetWithdrawals(context, service))
	r.Get("/ping", CheckConn(context, service))
	return r
}

func GetWithdrawals(ctx context.Context, service *gophermart.Gophermart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.GetWithdrawals(ctx)
	}
}

func CheckConn(ctx context.Context, service *gophermart.Gophermart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func WithdrawBalance(ctx context.Context, service *gophermart.Gophermart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var withdrawRequest Openapi.UserBalanceWithdrawRequest
		err := json.NewDecoder(r.Body).Decode(&withdrawRequest)
		if err != nil || withdrawRequest.Order == "" || withdrawRequest.Sum == 0 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		login := userLoginFromRequest(r)

		if err := service.WithDrawUserBalance(ctx, login, withdrawRequest); err != nil {
			if err == gophermart.ErrNoMoney {
				w.WriteHeader(http.StatusPaymentRequired)
				return
			}
			if err == gophermart.ErrOrderNotFound {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func GetUserBalance(ctx context.Context, service *gophermart.Gophermart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := userLoginFromRequest(r)
		body, err := service.GetUserBalance(ctx, login)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(body); err != nil {
			http.Error(w, "Unmarshalling error", http.StatusBadRequest)
			return
		}
	}
}

func GetUserOrders(ctx context.Context, service *gophermart.Gophermart) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userOrders := service.GetUserOrders(ctx, userLoginFromRequest(r))
		if len(userOrders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		var orders Openapi.GetUserOrdersResponse
		for _, order := range userOrders {
			orders = append(orders, Openapi.Order{Accrual: &order.Accrual, Number: order.OrderID, Status: order.Status, UploadedAt: order.UploadedAt.Format("2006-01-02 15:04:05Z07:00")})
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			http.Error(w, "Unmarshalling error", http.StatusBadRequest)
			return
		}
	}
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
