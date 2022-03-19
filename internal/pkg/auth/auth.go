package auth

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var JwtKey = []byte("diploma_prj_key") // можно получить откуда-нибудь

type UserClaims struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

func CreateToken(username string) (http.Cookie, error) {
	expiredAt := time.Now().Add(time.Hour * 24)
	claims := UserClaims{
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt.Unix(),
			Issuer:    "Gophermart",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(JwtKey)
	if err != nil {
		return http.Cookie{}, err
	}
	cookie := http.Cookie{
		Name:    "token",
		Value:   signedToken,
		Expires: expiredAt,
	}
	return cookie, err
}

func CheckToken(token string) (string, error) {
	claims := &UserClaims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		return "", err
	}
	if tkn.Valid {
		return claims.UserName, nil
	}
	return "", nil
}
