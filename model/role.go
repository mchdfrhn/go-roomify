// Package model meredefinisikan entitas data aplikasi.
package model

// Role merepresentasikan peran pengguna dalam sistem autentikasi dan otorisasi.
type Role struct {
	Id       string `json:"id"`       // ID unik role (UUID)
	Position string `json:"position"` // Nama posisi/jabatan peran (contoh: ADMIN, USER)
}
