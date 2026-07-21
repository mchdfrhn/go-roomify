// Package request menyimpan DTO permintaan client.
package request

// UserCredentialRequest merepresentasikan payload untuk proses registrasi atau login akun pengguna.
type UserCredentialRequest struct {
	Username string `json:"username" binding:"required"` // Username login pengguna (wajib)
	Password string `json:"password" binding:"required"` // Password login pengguna (wajib)
}

// UserUpdatePasswordRequest merepresentasikan payload permintaan pengubahan password pengguna.
type UserUpdatePasswordRequest struct {
	Id          string `json:"id" binding:"required"`           // ID kredensial pengguna (wajib)
	OldPassword string `json:"old_password" binding:"required"` // Password lama pengguna (wajib)
	NewPassword string `json:"new_password" binding:"required"` // Password baru pengguna (wajib)
}
