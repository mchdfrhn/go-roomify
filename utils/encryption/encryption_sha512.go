package encryption

import (
	"crypto/sha512"
	"encoding/hex"
)

func Sha512HashPassword(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func Sha512ComparePassword(password string, password_hash string) bool {

	return Sha512HashPassword(password) == password_hash
}
