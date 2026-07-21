// Package response menentukan bentuk pembungkusan (wrapper) standar balasan JSON API.
package response

// Import DTO pagination untuk embed struktur Paging.
import "go-roomify/model/dto" // Import DTO pagination

// Status merepresentasikan status eksekusi API (kode HTTP dan deskripsi pesan).
type Status struct {
	Code        int    `json:"code"`        // Kode status numerik HTTP
	Description string `json:"description"` // Pesan deskripsi status
}

// PagedResponse merepresentasikan standar response API untuk data list berhalaman (pagination).
type PagedResponse struct {
	Status Status        `json:"status"` // Header status response
	Data   []interface{} `json:"data"`   // Array data item yang dikembalikan
	Paging dto.Paging    `json:"paging"` // Metadata paginasi
}

// SingleResponse merepresentasikan standar response API untuk data tunggal (single resource).
type SingleResponse struct {
	Status Status      `json:"status"` // Header status response
	Data   interface{} `json:"data"`   // Object data payload utama
}
