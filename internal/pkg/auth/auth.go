package auth

import (
	"errors"
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

func GetClaims(jwtToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		jwtToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(JwtKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}
	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return nil, errors.New("jwt is expired")
	}
	return claims, nil
}
