// Package model meredefinisikan entitas data aplikasi.
package model

// ReservationStatus merepresentasikan status dari pengajuan reservasi ruangan.
type ReservationStatus struct {
	Id   string `json:"id"`   // ID unik status reservasi (UUID / ID acuan)
	Name string `json:"name"` // Nama status (contoh: PENDING, APPROVED, REJECTED, CANCELLED)
}
