// Package model meredefinisikan entitas data aplikasi.
package model

// UserProfile merepresentasikan data profil lengkap seorang pengguna aplikasi, termasuk relasi divisi, akun kredensial, dan perannya.
type UserProfile struct {
	Id          string         `json:"id" binding:"required"`              // ID unik profil pengguna (UUID)
	FullName    string         `json:"full_name" binding:"required"`       // Nama lengkap pengguna
	Division    Division       `json:"division" binding:"required"`        // Objek data divisi tempat pengguna bertugas
	Address     string         `json:"address" binding:"required"`         // Alamat domisili pengguna
	PhoneNumber string         `json:"phone_number" binding:"required"`    // Nomor telepon pengguna
	User        UserCredential `json:"user_credential" binding:"required"` // Objek data akun kredensial login pengguna
	Role        Role           `json:"role" binding:"required"`            // Objek data peran/akses pengguna dalam aplikasi
}
