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
func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &UserClaims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
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
