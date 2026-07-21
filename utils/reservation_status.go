// Package utils menyediakan nilai-nilai konstanta dan bantuan utilitas.
package utils

// Deklarasi konstanta status reservasi ruangan dalam bentuk string ID.
const (
	RESERV_STATUS_PENDING  = "1" // Status reservasi menunggu persetujuan (Pending)
	RESERV_STATUS_CANCEL   = "2" // Status reservasi dibatalkan oleh pengguna (Cancel)
	RESERV_STATUS_ACCEPTED = "3" // Status reservasi disetujui (Accepted/Approved)
	RESERV_STATUS_DECLINE  = "4" // Status reservasi ditolak oleh admin/GA (Decline)
)
