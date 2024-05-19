package adminExists

import (
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	cost := 10
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func passwordValid(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
