package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(hashedPwd string, plainPwd []byte) bool {
	byteHashedPwd := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHashedPwd, plainPwd)
	if err != nil {
		return false
	}
	return true
}
