// Package request menyimpan DTO permintaan client.
package request

// ReservationRequest merepresentasikan payload permintaan untuk pembuatan atau pengajuan reservasi ruangan.
type ReservationRequest struct {
	Id                 string                     `json:"id"`                   // ID reservasi (opsional saat pengajuan baru)
	UserProfileId      string                     `json:"user_profile_id"`      // ID profil pengguna pemesan ruangan
	ReservationDate    string                     `json:"reservation_date"`     // Tanggal reservasi diajukan
	StartDate          string                     `json:"start_date" binding:"required"` // Waktu mulai penggunaan ruangan (wajib)
	EndDate            string                     `json:"end_date" binding:"required"`   // Waktu selesai penggunaan ruangan (wajib)
	StatusId           string                     `json:"status_id"`            // ID status reservasi (contoh: Pending, Approved)
	RoomId             string                     `json:"room_id"  binding:"required"`   // ID ruangan yang akan dipesan (wajib)
	RequestMessage     string                     `json:"request_message"`      // Catatan/alasan pengajuan pemesanan dari user
	ResponseMessage    string                     `json:"response_message"`     // Catatan balasan dari admin
	AdditionalFacility []ReservationDetailRequest `json:"additional_facility"`  // Daftar fasilitas tambahan yang dipesan
}

// ReservationDetailRequest merepresentasikan rincian fasilitas tambahan dalam sebuah reservasi.
type ReservationDetailRequest struct {
	Id            string `json:"id"`             // ID detail reservasi (UUID)
	ReservationId string `json:"reservation_id"` // ID reservasi utama (foreign key)
	FacilityId    string `json:"facility_id"`    // ID fasilitas yang diminta (foreign key)
}

// ReservationStatusRequest merepresentasikan payload perubahan status reservasi oleh Admin/Approver.
type ReservationStatusRequest struct {
	ReservationId   string `json:"reservation_id" binding:"required"` // ID reservasi yang akan diubah statusnya (wajib)
	StatusId        string `json:"status_id"`                         // ID status baru (Approved / Rejected)
	ResponseMessage string `json:"response_message"`                  // Pesan balasan persetujuan/penolakan dari admin
}

// ReservationGetListFilter merepresentasikan kumpulan parameter filter untuk pencarian daftar reservasi.
type ReservationGetListFilter struct {
	UserId           string // Filter berdasarkan ID pengguna
	UserRole         string // Peran pengguna yang melakukan query
	ReservationId    string // Filter berdasarkan ID spesifik reservasi
	FilterStatus     string // Filter berdasarkan ID status reservasi
	FilterStartDate  string // Filter rentang waktu mulai reservasi
	FilterEndDate    string // Filter rentang waktu selesai reservasi
	FilterRoomId     string // Filter berdasarkan ID ruangan
	FilterPageNumber int    // Halaman pencarian saat ini
	FilterPageSize   int    // Jumlah data per halaman
}
