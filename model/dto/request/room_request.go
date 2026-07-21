// Package request menyimpan DTO permintaan client.
package request

// RoomRequest merepresentasikan payload pembuatan data ruangan baru.
type RoomRequest struct {
	Name          string `json:"name" binding:"required"`           // Nama ruangan (wajib)
	RoomTypeId    string `json:"room_type_id" binding:"required"`   // ID kategori tipe ruangan (wajib)
	Capacity      int    `json:"capacity" binding:"required"`       // Kapasitas penampungan orang (wajib)
	IsReserveable *bool  `json:"is_reserveable" binding:"required"` // Status dapat dipesan atau tidak (wajib)
}

// RoomStatusRequest merepresentasikan payload pengubahan status ketersediaan ruangan.
type RoomStatusRequest struct {
	Id          string `json:"id" binding:"required"`           // ID ruangan (wajib)
	IsAvailable *bool  `json:"is_available" binding:"required"` // Status ketersediaan ruangan saat ini (wajib)
}

// UpdateRoomRequest merepresentasikan payload pembaruan data ruangan secara keseluruhan.
type UpdateRoomRequest struct {
	Id            string `json:"id" binding:"required"`             // ID ruangan yang hendak diperbarui (wajib)
	Name          string `json:"name" binding:"required"`           // Nama ruangan baru (wajib)
	RoomTypeId    string `json:"room_type_id" binding:"required"`   // ID kategori tipe ruangan baru (wajib)
	Capacity      int    `json:"capacity" binding:"required"`       // Kapasitas penampungan baru (wajib)
	IsReserveable *bool  `json:"is_reserveable" binding:"required"` // Status dapat dipesan baru (wajib)
}
