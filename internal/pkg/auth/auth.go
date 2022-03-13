package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("diploma_prj_key") // можно получить откуда-нибудь

type UserClaims struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func createToken(username string) (string, error) {
	claims := UserClaims{
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "Gophermart",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func getClaims(jwtToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		jwtToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
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
