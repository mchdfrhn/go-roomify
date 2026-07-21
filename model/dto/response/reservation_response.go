// Package response menentukan bentuk pembungkusan standar balasan JSON API.
package response

// Import model entitas bisnis utama.
import "go-roomify/model" // Import model utama

// ReservationResponse merepresentasikan struktur lengkap data reservasi yang dikembalikan ke client.
type ReservationResponse struct {
	ID                 string                         `json:"id"`                  // ID unik reservasi (UUID)
	UserProfile        ReservationUserProfileResponse `json:"user_profile"`        // Informasi profil pemesan
	ReservationDate    string                         `json:"reservation_date"`    // Tanggal pengajuan pemesanan
	StartTime          string                         `json:"start_time"`          // Waktu mulai pemakaian
	EndTime            string                         `json:"end_time"`            // Waktu selesai pemakaian
	ReservationStatus  model.ReservationStatus        `json:"reservation_status"`  // Objek status reservasi saat ini
	Room               ReservationRoomResponse        `json:"room"`                // Detail ruangan yang dipesan
	RequestMessage     string                         `json:"request_message"`     // Catatan/alasan pengajuan pemesan
	ResponseMessage    string                         `json:"response_message"`    // Catatan tanggapan dari administrator
	AdditionalFacility []model.Facility               `json:"additional_facility"` // Daftar fasilitas tambahan yang dipesan
}

// ReservationUserProfileResponse merepresentasikan rincian profil singkat pemesan ruangan pada respon reservasi.
type ReservationUserProfileResponse struct {
	ID          string         `json:"id"`           // ID profil pemesan
	FullName    string         `json:"full_name"`    // Nama lengkap pemesan
	Division    model.Division `json:"division"`     // Objek divisi pemesan
	Role        model.Role     `json:"role"`         // Objek peran pemesan
	Address     string         `json:"address"`      // Alamat pemesan
	PhoneNumber string         `json:"phone_number"` // Nomor kontak pemesan
}

// ReservationRoomResponse merepresentasikan rincian ruangan yang dipesan pada respon reservasi.
type ReservationRoomResponse struct {
	ID            string           `json:"id"`             // ID unik ruangan
	Name          string           `json:"name"`           // Nama ruangan
	RoomType      model.RoomType   `json:"roomtype"`       // Objek kategori tipe ruangan
	Capacity      int              `json:"capacity"`       // Kapasitas ruangan
	IsAvailable   bool             `json:"is_available"`   // Status ketersediaan fisik
	IsReserveable bool             `json:"is_reserveable"` // Status reservabilitas
	Facility      []model.Facility `json:"facility"`       // Daftar fasilitas bawaan dalam ruangan
}
