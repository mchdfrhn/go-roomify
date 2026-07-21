// Package encryption menyediakan utilitas enkripsi.
package encryption

// Import standar library SHA-512 dan konversi heksadesimal.
import (
	"crypto/sha512" // Package enkripsi SHA-512
	"encoding/hex"   // Package enkoding string hex
)

// Sha512HashPassword menghasilkan string hash heksadesimal SHA-512 dari password.
func Sha512HashPassword(password string) string {
	// Inisialisasi hash SHA-512 baru
	hash := sha512.New()
	// Menginputkan byte string password ke dalam engine hash
	hash.Write([]byte(password))
	// Mengembalikan string heksadesimal dari checksum hash yang dihasilkan
	return hex.EncodeToString(hash.Sum(nil))
}

// Sha512ComparePassword membandingkan password plain text dengan hash SHA-512.
func Sha512ComparePassword(password string, password_hash string) bool {
	// Mengembalikan true jika hasil SHA-512 hash password sama dengan password_hash
	return Sha512HashPassword(password) == password_hash
}
