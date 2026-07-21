// Package repository menangani akses dan manipulasi data langsung ke database PostgreSQL.
package repository

// Import package database/sql, json, dto, utils, validation, strconv, query, dan errors.
import (
	"database/sql"                    // Interface koneksi database SQL
	"encoding/json"                   // Serialisasi/deserialisasi JSON
	"errors"                          // Membuat error instan
	"go-roomify/model/dto"            // DTO paginasi
	"go-roomify/model/dto/request"    // DTO request reservasi
	"go-roomify/model/dto/response"   // DTO response reservasi
	"go-roomify/utils"                // Helper konstanta & paginasi
	"go-roomify/utils/query"          // Builder query SQL dinamis
	"go-roomify/utils/validation"     // Helper escape string
	"strconv"                         // Konversi angka ke string
)

// ReservationRepository merupakan kontrak interface untuk mengelola transaksi pemesanan ruangan.
type ReservationRepository interface {
	CreateRequest(new_request request.ReservationRequest) error                                                      // Membuat permintaan reservasi baru
	ChangeStatus(reserv_status request.ReservationStatusRequest) error                                               // Mengubah status persetujuan reservasi
	GetListByToken(fl_reserv_get_list request.ReservationGetListFilter) ([]response.ReservationResponse, dto.Paging, error) // Mengambil daftar reservasi berdasar filter & token
}

// reservationRepository merupakan struktur konkrit pengelola transaksi reservasi.
type reservationRepository struct {
	db *sql.DB // Pointer koneksi database PostgreSQL
}

// CreateRequest menyisipkan data reservasi header dan fasilitas tambahan ke database.
func (self *reservationRepository) CreateRequest(new_request request.ReservationRequest) error {
	// Inisialisasi builder INSERT untuk tx_reservation
	qinsert := query.QInsert{DB: self.db}

	qinsert.Table("tx_reservation")
	qinsert.Column(
		"id",
		"user_profile_id",
		"reservation_date",
		"start_time",
		"end_time",
		"reservation_status_id",
		"room_id",
		"request_message")
	// Parameter nilai header reservasi
	qinsert.Values(
		new_request.Id,
		new_request.UserProfileId,
		new_request.ReservationDate,
		new_request.StartDate,
		new_request.EndDate,
		new_request.StatusId,
		new_request.RoomId,
		new_request.RequestMessage)

	// Eksekusi insert header
	result, err := qinsert.Run()

	if err != nil {
		return err
	}

	// Cek jumlah baris yang terpengaruh
	if a, _ := result.RowsAffected(); a == 0 {
		return errors.New("Failed to make new Room Reservation Request")
	}

	// Inisialisasi builder INSERT untuk detail fasilitas tambahan
	qinsert = query.QInsert{DB: self.db}

	qinsert.Table("tx_reservation_detail")
	qinsert.Column(
		"id",
		"reservation_id",
		"facility_id")

	// Iterasi fasilitas tambahan yang dipesan
	for _, additional_facility := range new_request.AdditionalFacility {
		qinsert.Values(
			additional_facility.Id,
			new_request.Id,
			additional_facility.FacilityId)
	}

	// Eksekusi insert detail fasilitas
	result, err = qinsert.Run()

	if err != nil {
		return err
	}

	if a, _ := result.RowsAffected(); a == 0 {
		return errors.New("Failed to make new Room Reservation Request")
	}

	return nil
}

// ChangeStatus memperbarui status reservasi (Disetujui/Ditolak) dan ketersediaan fisik ruangan.
func (self *reservationRepository) ChangeStatus(reserv_status request.ReservationStatusRequest) error {
	// Ambil room_id dari data reservasi yang diproses
	qselect := query.QSelect{DB: self.db}
	qselect.Table("tx_reservation AS rv")
	qselect.Column("rv.room_id")
	qselect.Where("rv.id", "=", reserv_status.ReservationId)

	var room_id string

	// Eksekusi pencarian room_id
	if err := qselect.RunRow().Scan(&room_id); err != nil {
		return err
	}

	// Jika status disetujui (ACCEPTED), ubah status fisik ketersediaan ruangan menjadi false
	if reserv_status.StatusId == utils.RESERV_STATUS_ACCEPTED {
		qupdate := query.QUpdate{DB: self.db}
		qupdate.Table("mst_room")
		qupdate.Set("is_available", false)
		qupdate.Where("id", "=", room_id)

		result, err := qupdate.Run()

		if err != nil {
			return err
		}

		if a, _ := result.RowsAffected(); a == 0 {
			return errors.New("Invalid Room Id")
		}
	}

	// Perbarui status transaksi reservasi dan pesan balasan dari GA/Admin
	qupdate := query.QUpdate{DB: self.db}

	qupdate.Table("tx_reservation")
	qupdate.Set("reservation_status_id", reserv_status.StatusId)
	qupdate.Set("response_message", reserv_status.ResponseMessage)
	qupdate.Where("id", "=", reserv_status.ReservationId)
	qupdate.AndWhere("reservation_status_id", "=", utils.RESERV_STATUS_PENDING)

	result, err := qupdate.Run()

	if err != nil {
		return err
	}

	if a, _ := result.RowsAffected(); a == 0 {
		return errors.New("Invalid Reservation Id")
	}

	return err
}

// GetListByToken mengambil daftar reservasi terstruktur JSON berdasarkan filter dan hak akses pengguna.
func (self *reservationRepository) GetListByToken(fl_reserv_get_list request.ReservationGetListFilter) ([]response.ReservationResponse, dto.Paging, error) {
	// Ekstraksi nilai filter
	filter_status := fl_reserv_get_list.FilterStatus
	filter_start_date := fl_reserv_get_list.FilterStartDate
	filter_end_date := fl_reserv_get_list.FilterEndDate
	filter_room_id := fl_reserv_get_list.FilterRoomId
	filter_pagenum := fl_reserv_get_list.FilterPageNumber
	filter_pagesize := fl_reserv_get_list.FilterPageSize
	resrv_id := fl_reserv_get_list.ReservationId
	user_id := fl_reserv_get_list.UserId
	user_role := fl_reserv_get_list.UserRole

	resultPagingDto := dto.Paging{}

	// Menyusun query SQL kompleks dengan subquery JSON_AGG & JSON_BUILD_OBJECT PostgreSQL
	var querySQL = `
	SELECT JSON_AGG(reservation_info) AS all_reservations
	FROM (
		SELECT JSON_BUILD_OBJECT(
			'id', tx_reservation.id,
			'user_profile', (
				SELECT JSON_BUILD_OBJECT(
					'id', mst_user_profile.id,
					'full_name', mst_user_profile.full_name,
					'division', JSON_BUILD_OBJECT(
						'id', mst_division.id,
						'name', mst_division.name
					),
					'role', JSON_BUILD_OBJECT(
						'id', mst_role.id,
						'position', mst_role.position
					),
					'address', mst_user_profile.address,
					'phone_number', mst_user_profile.phone_number
				)
				FROM mst_user_profile
				LEFT JOIN mst_division ON mst_user_profile.division_id = mst_division.id
				LEFT JOIN mst_role ON mst_user_profile.role_id = mst_role.id
				WHERE mst_user_profile.id = tx_reservation.user_profile_id
			),
			'reservation_date', tx_reservation.reservation_date,
			'start_time', tx_reservation.start_time,
			'end_time', tx_reservation.end_time,
			'reservation_status', (
				SELECT JSON_BUILD_OBJECT(
					'id', tx_reservation_status.id,
					'name', tx_reservation_status.name
				)
				FROM tx_reservation_status
				WHERE tx_reservation_status.id = tx_reservation.reservation_status_id
			),
			'room', (
				SELECT JSON_BUILD_OBJECT(
					'id', mst_room.id,
					'name', mst_room.name,
					'roomtype', JSON_BUILD_OBJECT(
						'id', room_type.id,
						'name', room_type.name
					),
					'capacity', mst_room.capacity,
					'is_available', mst_room.is_available,
					'is_reserveable', mst_room.is_reserveable,
					'facility', (
						SELECT COALESCE(
							JSON_AGG(
								JSON_BUILD_OBJECT(
									'id', mst_facility.id,
									'name', mst_facility.name,
									'is_available', mst_facility.is_available,
									'is_reserveable', mst_facility.is_reserveable,
									'room_id', mst_facility.room_id
								)
							), '[]'::JSON
						)
						FROM mst_facility
						WHERE mst_facility.room_id = mst_room.id
					)
				)
				FROM mst_room
				JOIN room_type
					ON room_type.id = mst_room.room_type_id
				WHERE mst_room.id = tx_reservation.room_id
			),
			'request_message', tx_reservation.request_message,
			'response_message', tx_reservation.response_message,
			'additional_facility', (
				SELECT COALESCE(
					JSON_AGG(
						JSON_BUILD_OBJECT(
							'id', mst_facility.id,
							'name', mst_facility.name,
							'is_available', mst_facility.is_available,
							'is_reserveable', mst_facility.is_reserveable,
							'room_id', mst_facility.room_id
						)
					), '[]'::JSON
				)
				FROM tx_reservation_detail
				LEFT JOIN mst_facility ON tx_reservation_detail.facility_id = mst_facility.id
				WHERE tx_reservation_detail.reservation_id = tx_reservation.id
			)
		) AS reservation_info
		FROM tx_reservation
		LEFT JOIN tx_reservation_status ON tx_reservation_status.id = tx_reservation.reservation_status_id
		LEFT JOIN mst_room ON mst_room.id = tx_reservation.room_id
		WHERE true `
	// Jika role adalah employee, filter hanya untuk id miliknya sendiri
	if user_role == "employee" {
		querySQL += "AND tx_reservation.user_profile_id = " + validation.EscapeString(user_id)
	}

	// Filter berdasarkan ID spesifik atau kombinasi filter lain
	if resrv_id != "" {
		querySQL += "AND tx_reservation.id = " + validation.EscapeString(resrv_id)
	} else {
		if filter_status != "" {
			querySQL += "AND tx_reservation_status.id = " + validation.EscapeString(filter_status)
		}
		if filter_start_date != "" {
			querySQL += "AND tx_reservation.reservation_date >= " + validation.EscapeString(filter_start_date)
		}
		if filter_end_date != "" {
			querySQL += "AND tx_reservation.reservation_date <= " + validation.EscapeString(filter_end_date)
		}
		if filter_room_id != "" {
			querySQL += "AND mst_room.id = " + validation.EscapeString(filter_room_id)
		}
	}
	// Menghitung offset paginasi
	skip := (filter_pagenum - 1) * filter_pagesize

	// Menambahkan LIMIT dan OFFSET pada query
	querySQL += " LIMIT " + strconv.Itoa(filter_pagesize)
	querySQL += " OFFSET " + strconv.Itoa(skip)

	querySQL += `
	) AS all_reservations`

	var bjson_reserv []byte
	var rows_reserv_response []response.ReservationResponse

	// Eksekusi query untuk menerima raw JSON byte
	err := self.db.QueryRow(querySQL).Scan(&bjson_reserv)

	if err != nil {
		return nil, resultPagingDto, err
	}

	// Memeriksa jika hasil query kosong
	if len(bjson_reserv) == 0 {
		return rows_reserv_response, resultPagingDto, nil
	}

	// Unmarshal byte JSON ke slice struct response.ReservationResponse
	err = json.Unmarshal(bjson_reserv, &rows_reserv_response)

	if err != nil {
		return nil, resultPagingDto, errors.New("Error unmarshaling JSON")
	}

	// Menghitung total data reservasi untuk metadata paginasi
	var totalRows int
	qcount := query.QSelect{DB: self.db}
	qcount.Table("tx_reservation")
	qcount.Column("COUNT(id)")

	count_err := qcount.RunRow().Scan(&totalRows)
	if err != nil {
		return nil, resultPagingDto, count_err
	}

	// Mengolah DTO metadata paginasi
	resultPagingDto = utils.Paginate(filter_pagenum, filter_pagesize, totalRows)

	return rows_reserv_response, resultPagingDto, nil
}

// NewReservationRepository menginisialisasi instansi konkrit ReservationRepository baru.
func NewReservationRepository(db *sql.DB) ReservationRepository {
	return &reservationRepository{
		db: db,
	}
}
