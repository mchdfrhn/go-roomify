package encryption

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
	if err != nil {
		return "", err
	}

	return string(password_hash), nil
}

func ComparePassword(password string, password_hash string) (bool) {
	err := bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))
	
	return err == nil
}
