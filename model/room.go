// Package model meredefinisikan entitas data aplikasi.
package model

// Room merepresentasikan entitas ruangan fisik yang tersedia untuk dipesan.
type Room struct {
	Id            string `json:"id" binding:"required"`             // ID unik ruangan (UUID)
	Name          string `json:"name" binding:"required"`           // Nama ruangan (contoh: Ruang Rapat Lt. 2)
	RoomTypeId    string `json:"room_type_id" binding:"required"`   // ID kategori tipe ruangan (foreign key ke RoomType)
	Capacity      int    `json:"capacity" binding:"required"`       // Kapasitas maksimum jumlah orang
	IsAvailable   *bool  `json:"is_available" binding:"required"`   // Pointer ke bool status ketersediaan (bebas/dipakai)
	IsReserveable *bool  `json:"is_reserveable" binding:"required"` // Pointer ke bool status apakah ruangan bisa dipesan
}
