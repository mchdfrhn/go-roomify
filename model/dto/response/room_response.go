// Package response menentukan bentuk pembungkusan standar balasan JSON API.
package response

// FacilityForRoomResponse merepresentasikan rincian fasilitas singkat di dalam objek respon ruangan.
type FacilityForRoomResponse struct {
	Id     string `json:"id"`      // ID fasilitas (UUID)
	Name   string `json:"name"`    // Nama fasilitas
	RoomId string `json:"room_id"` // ID ruangan pemiliki fasilitas
}

// RoomResponse merepresentasikan DTO informasi data ruangan secara mendalam beserta fasilitasnya.
type RoomResponse struct {
	Id            string                    `json:"id"`             // ID unik ruangan (UUID)
	Name          string                    `json:"name"`           // Nama ruangan
	RoomType      string                    `json:"roomtype"`       // Nama/deskripsi tipe ruangan
	Capacity      int                       `json:"capacity"`       // Kapasitas penampungan
	IsAvailable   bool                      `json:"is_available"`   // Status ketersediaan fisik
	StartTime     string                    `json:"start_time"`     // Waktu awal jam operasional/reservasi
	EndTIme       string                    `json:"end_time"`       // Waktu akhir jam operasional/reservasi
	IsReserveable bool                      `json:"is_reserveable"` // Status reservabilitas
	Facilities    []FacilityForRoomResponse `json:"facilities"`     // Slice fasilitas yang berada di dalam ruangan
}
