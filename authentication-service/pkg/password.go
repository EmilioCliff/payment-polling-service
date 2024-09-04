package pkg

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string, cost int) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("Couldn't hash password: %s", err)
	}

	return string(hashPassword), nil
}

func ComparePasswordAndHash(hashPass string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
}
