// Package model meredefinisikan entitas data aplikasi.
package model

// UserCredentialJwt merepresentasikan data kredensial pengguna beserta role dan token JWT yang dihasilkan saat login/autentikasi.
type UserCredentialJwt struct {
	Id       string `json:"id"`       // ID unik kredensial pengguna (UUID)
	Username string `json:"username" binding:"required"` // Username login
	Password string `json:"password" binding:"required"` // Kata sandi (terenkripsi)
	Role     string `json:"role"`     // Nama role/posisi pengguna
	Token    string `json:"token"`    // Token bearer JWT yang aktif
}

// UserCredential merepresentasikan data kredensial akun pengguna tanpa role tambahan.
type UserCredential struct {
	Id       string `json:"id"`       // ID unik kredensial (UUID)
	Username string `json:"username" binding:"required"` // Username login
	Password string `json:"password" binding:"required"` // Kata sandi (terenkripsi)
	Token    string `json:"token"`    // Token pengenal/autentikasi
}
