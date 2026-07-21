// Package response menentukan bentuk pembungkusan standar balasan JSON API.
package response

// UserCredentialResponse merepresentasikan payload respon autentikasi sukses (login) berisi token JWT dan User ID.
type UserCredentialResponse struct {
	AccessToken string `json:"access_token"` // Token JWT untuk otorisasi endpoint terproteksi
	UserId      string `json:"user_id"`      // ID kredensial pengguna yang berhasil login
}
