package controller

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go-musthave-diploma-tpl/internal/pkg/auth"
	"net/http"
	"strings"
)

type CtxKey string

const UserCtxKey = CtxKey("UserID")

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "register") && !strings.Contains(r.URL.Path, "login") {
			c, err := r.Cookie("token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			tknStr := c.Value
			userLogin, err := auth.CheckToken(tknStr)
			if err != nil {
				if errors.Is(err, jwt.ErrSignatureInvalid) {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if userLogin == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserCtxKey, userLogin)))
		}
		next.ServeHTTP(w, r)
	})
}
