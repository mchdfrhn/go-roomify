// Package encryption menyediakan utilitas enkripsi password berbasis bcrypt dan hash.
package encryption

// Import library golang bcrypt untuk enkripsi password aman.
import (
	"golang.org/x/crypto/bcrypt" // Library enkripsi bcrypt
)

// HashPassword membuat bcrypt hash dari plain text password.
func HashPassword(password string) (string, error) {
	// Meng-generate hash password dengan kostumisasi default cost bcrypt
	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Memeriksa apakah terdapat kesalahan saat pembuatan hash
	if err != nil {
		// Mengembalikan string kosong dan error jika gagal hash
		return "", err
	}

	// Mengembalikan string hash password dan nil error
	return string(password_hash), nil
}

// ComparePassword membandingkan plain text password dengan hash bcrypt.
func ComparePassword(password string, password_hash string) bool {
	// Membandingkan hash password dengan plain text password
	err := bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password))

	// Mengembalikan true jika cocok (err == nil), dan false jika salah
	return err == nil
}
