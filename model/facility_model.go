// Package model meredefinisikan entitas data aplikasi.
package model

// Facility merepresentasikan fasilitas pendukung yang terdapat di dalam sebuah ruangan.
type Facility struct {
	Id            string `json:"id"`             // ID unik fasilitas (UUID)
	Name          string `json:"name"`           // Nama fasilitas (contoh: Proyektor, AC, Whiteboard)
	IsAvailable   bool   `json:"is_available"`   // Status ketersediaan kondisi fisik fasilitas (baik/rusak)
	IsReserveable bool   `json:"is_reserveable"` // Status apakah fasilitas dapat dipinjam secara parsial
	RoomId        string `json:"room_id"`        // ID ruangan tempat fasilitas berada (foreign key)
}
