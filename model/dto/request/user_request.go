// Package request menyimpan DTO permintaan client.
package request

// UserProfileRequest merepresentasikan payload permintaan pembuatan atau pembaruan profil lengkap pengguna.
type UserProfileRequest struct {
	Id          string `json:"id"`                         // ID profil (opsional saat pendaftaran)
	FullName    string `json:"full_name" binding:"required"`    // Nama lengkap pengguna (wajib)
	DivisionId  string `json:"division_id" binding:"required"`  // ID divisi asal pengguna (wajib)
	Address     string `json:"address" binding:"required"`      // Alamat domisili pengguna (wajib)
	PhoneNumber string `json:"phone_number" binding:"required"` // Nomor telepon kontak (wajib)
	UserId      string `json:"user_id" binding:"required"`      // ID akun kredensial terasosiasi (wajib)
	RoleId      string `json:"role_id" binding:"required"`      // ID peran pengguna (wajib)
}
