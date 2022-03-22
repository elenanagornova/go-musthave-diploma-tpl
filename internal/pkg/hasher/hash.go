package hasher

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(hashedPwd string, plainPwd []byte) bool {
	byteHashedPwd := []byte(hashedPwd)
	return bcrypt.CompareHashAndPassword(byteHashedPwd, plainPwd) == nil
}
