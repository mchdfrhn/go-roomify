// Package request menyimpan DTO permintaan client.
package request

// RoomTypeRequest merepresentasikan payload pembuatan atau pembaruan nama kategori/tipe ruangan.
type RoomTypeRequest struct {
	Name string `json:"name" binding:"required"` // Nama tipe ruangan (wajib)
}
