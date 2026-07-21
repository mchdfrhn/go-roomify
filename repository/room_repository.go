// Package repository menangani akses dan manipulasi data langsung ke database PostgreSQL.
package repository

// Import package database/sql, fmt, model, dto, request/response, utils, dan query builder.
import (
	"database/sql"                    // Interface koneksi database SQL
	"fmt"                             // Formatting string
	"go-roomify/model"                // Struct model data Ruangan
	"go-roomify/model/dto"            // DTO paginasi
	"go-roomify/model/dto/request"    // DTO request perbaruan/filter ruangan
	"go-roomify/model/dto/response"   // DTO response ruangan
	"go-roomify/utils"                // Helper paginasi
	"go-roomify/utils/query"          // Builder query SQL dinamis
)

// RoomRepository merupakan kontrak interface untuk mengelola data master dan status ruangan.
type RoomRepository interface {
	CreateRoom(roomModel model.Room) error                                                                     // Menyimpan record ruangan baru
	GetRoomIdIfExist(name string, roomtypeId string) (string, error)                                            // Cek keberadaan ID ruangan berdasarkan nama dan tipe
	GetAllRoom(page int, skip int, size int, paramType string) ([]response.RoomResponse, dto.Paging, error)     // Mengambil semua ruangan berhalaman
	GetRoomByIdOrName(idOrNameRoom string) ([]response.RoomResponse, error)                                     // Mencari ruangan berdasarkan ID atau nama
	UpdateRoomById(updateRoom request.UpdateRoomRequest) (request.UpdateRoomRequest, error)                    // Memperbarui atribut utama ruangan
	DeleteRoomById(roomId string) error                                                                        // Menghapus ruangan berdasarkan ID
	GetAvailableRoom(paramType string) ([]response.RoomResponse, error)                                         // Mengambil daftar ruangan yang sedang tersedia
	UpdateRoomByIdAvailableOnly(updateRoom request.RoomStatusRequest) error                                    // Memperbarui status ketersediaan (is_available) ruangan saja
}

// roomRepository merupakan struktur konkrit pengelola data ruangan.
type roomRepository struct {
	db *sql.DB // Pointer koneksi database PostgreSQL
}

// CreateRoom menyisipkan data ruangan baru ke tabel mst_room.
func (rr *roomRepository) CreateRoom(roomModel model.Room) error {
	// Inisialisasi builder INSERT
	query := query.QInsert{DB: rr.db}

	// Menyiapkan kolom dan nilai ruangan
	_, err := query.Table(
		"mst_room",
	).Column(
		"id", "name", "room_type_id", "capacity", "is_available", "is_reserveable",
	).Values(
		roomModel.Id, roomModel.Name, roomModel.RoomTypeId, roomModel.Capacity,
		roomModel.IsAvailable, roomModel.IsReserveable,
	).Run()

	if err != nil {
		return err
	}

	return nil
}

// GetRoomIdIfExist mencari ID ruangan yang cocok dengan kombinasi nama dan ID tipe ruangan.
func (rr *roomRepository) GetRoomIdIfExist(name string, roomtypeId string) (string, error) {
	// Inisialisasi builder SELECT
	query := query.QSelect{DB: rr.db}
	var idRoom string

	// Eksekusi pencarian dengan dua kriteria WHERE
	err := query.Table(
		"mst_room",
	).Column(
		"id",
	).Where(
		"name", "=", name,
	).AndWhere(
		"room_type_id", "=", roomtypeId,
	).RunRow().Scan(
		&idRoom,
	)

	if err != nil {
		// Mengabaikan error jika memang tidak ada baris yang cocok (ErrNoRows)
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return idRoom, nil
}

// GetAllRoom mengambil seluruh data ruangan beserta fasilitas dan reservasi terakhirnya secara paginasi.
func (rr *roomRepository) GetAllRoom(page int, skip int, size int, paramType string) ([]response.RoomResponse, dto.Paging, error) {
	queryAllRoom := query.QSelect{DB: rr.db}

	// Subquery tabel untuk membatasi baris ruangan sebelum di-JOIN dengan fasilitas
	subQueryTable := fmt.Sprintf(`
		(
			SELECT r.*, res.start_time, res.end_time
			FROM mst_room AS r
			LEFT JOIN (
				SELECT room_id, start_time, end_time
				FROM tx_reservation
				ORDER BY start_time DESC
				LIMIT 1
			) AS res ON r.id = res.room_id
			LIMIT %d OFFSET %d
		) AS r
	`, size, skip)

	// Menyusun relasi JOIN dengan tabel room_type dan mst_facility
	getRoom := queryAllRoom.Table(
		subQueryTable,
	).Column(
		"r.id",
		"r.name",
		"rt.name AS room_type",
		"r.capacity",
		"r.is_available",
		"COALESCE(r.start_time, '1945-08-17')",
		"COALESCE(r.end_time, '1945-08-17')",
		"r.is_reserveable",
		"COALESCE(f.id, 'null')",
		"COALESCE(f.name, 'null')",
		"COALESCE(f.room_id, 'null')",
	).LeftJoin(
		"mst_facility AS f", "f.room_id = r.id",
	).Join(
		"room_type AS rt", "rt.id = r.room_type_id",
	)

	// Filter berdasarkan tipe ruangan jika ada
	if paramType != "" {
		getRoom.Where(
			"rt.name", "=", paramType,
		)
	}

	// Urutkan berdasarkan kapasitas ruangan secara ascending
	rows, err := getRoom.OrderBy(
		"r.capacity", "ASC",
	).Run()

	if err != nil {
		return nil, dto.Paging{}, err
	}

	// Grouping data fasilitas ke dalam tiap struct RoomResponse
	responseData, err := rr.scanRoomAndFacility(rows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	// Hitung total baris ruangan
	var totalRows int
	queryCount := query.QSelect{DB: rr.db}
	err = queryCount.Table(
		"mst_room",
	).Column(
		"COUNT(id)",
	).RunRow().Scan(&totalRows)

	if err != nil {
		return nil, dto.Paging{}, err
	}

	// Olah DTO paginasi
	resultPagingDto := utils.Paginate(page, size, totalRows)
	return responseData, resultPagingDto, nil
}

// GetRoomByIdOrName mencari ruangan berdasarkan pencocokan ID atau Nama Ruangan.
func (rr *roomRepository) GetRoomByIdOrName(roomidOrName string) ([]response.RoomResponse, error) {
	query := query.QSelect{DB: rr.db}

	// Query pencarian gabungan dengan OR WHERE
	rows, err := query.Table(
		"mst_room AS r",
	).Column(
		"r.id",
		"r.name",
		"rt.name AS room_type",
		"r.capacity",
		"r.is_available",
		"COALESCE(res.start_time, '1945-08-17')",
		"COALESCE(res.end_time, '1945-08-17')",
		"r.is_reserveable",
		"COALESCE(f.id, 'null')",
		"COALESCE(f.name, 'null')",
		"COALESCE(f.room_id, 'null')",
	).LeftJoin(
		"mst_facility AS f", "f.room_id = r.id",
	).Join(
		"room_type As rt", "rt.id = r.room_type_id",
	).LeftJoin(
		`
		( 
			SELECT room_id, start_time, end_time
			FROM tx_reservation
			ORDER BY start_time DESC
			LIMIT 1 
		) AS res
		`,
		"r.id = res.room_id",
	).Where(
		"r.id", "=", roomidOrName,
	).OrWhere(
		"r.name", "=", roomidOrName,
	).OrderBy(
		"r.capacity", "ASC",
	).Run()

	if err != nil {
		return nil, err
	}

	// Scan data ke struct DTO
	responseData, err := rr.scanRoomAndFacility(rows)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

// UpdateRoomById memperbarui detail atribut ruangan (Nama, Tipe, Kapasitas, Dapat Direservasi).
func (rr *roomRepository) UpdateRoomById(updateRoom request.UpdateRoomRequest) (request.UpdateRoomRequest, error) {
	query := query.QUpdate{DB: rr.db}

	// Eksekusi perbaruan data
	_, err := query.Table(
		"mst_room",
	).Set(
		"name",
		updateRoom.Name,
	).Set(
		"room_type_id",
		updateRoom.RoomTypeId,
	).Set(
		"capacity",
		updateRoom.Capacity,
	).Set(
		"is_reserveable",
		updateRoom.IsReserveable,
	).Where(
		"id", "=", updateRoom.Id,
	).Run()

	if err != nil {
		return request.UpdateRoomRequest{}, err
	}

	return updateRoom, nil
}

// DeleteRoomById menghapus record ruangan dari database.
func (rr *roomRepository) DeleteRoomById(roomId string) error {
	query := query.QDelete{DB: rr.db}

	// Eksekusi query hapus
	_, err := query.Table(
		"mst_room",
	).Where(
		"id", "=", roomId,
	).Run()

	if err != nil {
		return err
	}

	return nil
}

// GetAvailableRoom mengambil hanya ruangan yang flag is_available = true.
func (rr *roomRepository) GetAvailableRoom(paramType string) ([]response.RoomResponse, error) {
	query := query.QSelect{DB: rr.db}

	// Menyusun query khusus ketersediaan
	getRoom := query.Table(
		"mst_room AS r",
	).Column(
		"r.id",
		"r.name",
		"rt.name AS room_type",
		"r.capacity",
		"r.is_available",
		"COALESCE(res.start_time, '1945-08-17')",
		"COALESCE(res.end_time, '1945-08-17')",
		"r.is_reserveable",
		"COALESCE(f.id, 'null')",
		"COALESCE(f.name, 'null')",
		"COALESCE(f.room_id, 'null')",
	).LeftJoin(
		"mst_facility AS f", "f.room_id = r.id",
	).Join(
		"room_type AS rt", "rt.id = r.room_type_id",
	).LeftJoin(
		`
		( 
			SELECT room_id, start_time, end_time
			FROM tx_reservation
			ORDER BY start_time DESC
			LIMIT 1 
		) AS res
		`,
		"r.id = res.room_id",
	).Where(
		"r.is_available", "=", true,
	)

	if paramType != "" {
		getRoom.AndWhere(
			"rt.name", "=", paramType,
		)
	}

	rows, err := getRoom.OrderBy(
		"r.capacity", "ASC",
	).Run()

	if err != nil {
		return nil, err
	}

	// Scanning data
	responseData, err := rr.scanRoomAndFacility(rows)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

// UpdateRoomByIdAvailableOnly memperbarui kolom status is_available saja pada suatu ruangan.
func (rr *roomRepository) UpdateRoomByIdAvailableOnly(updateRoom request.RoomStatusRequest) error {
	query := query.QUpdate{DB: rr.db}

	// Eksekusi update ketersediaan
	_, err := query.Table(
		"mst_room",
	).Set(
		"is_available",
		updateRoom.IsAvailable,
	).Where(
		"id", "=", updateRoom.Id,
	).Run()

	if err != nil {
		return err
	}

	return nil
}

// scanRoomAndFacility helper internal untuk mengelompokkan baris JOIN fasilitas menjadi slice RoomResponse terstruktur.
func (rr *roomRepository) scanRoomAndFacility(rows *sql.Rows) ([]response.RoomResponse, error) {
	var responseData []response.RoomResponse
	var roomResponse response.RoomResponse

	// Iterasi hasil baris database
	for rows.Next() {
		var currentRoom model.Room
		var currentFacility response.FacilityForRoomResponse
		var startTime string
		var endTime string

		// Scan baris data
		err := rows.Scan(
			&currentRoom.Id,
			&currentRoom.Name,
			&currentRoom.RoomTypeId,
			&currentRoom.Capacity,
			&currentRoom.IsAvailable,
			&startTime,
			&endTime,
			&currentRoom.IsReserveable,
			&currentFacility.Id,
			&currentFacility.Name,
			&currentFacility.RoomId,
		)
		if err != nil {
			return nil, err
		}

		// Jika ID ruangan berubah, simpan roomResponse sebelumnya dan buat objek roomResponse baru
		if currentRoom.Id != roomResponse.Id {
			if roomResponse.Id != "" {
				responseData = append(responseData, roomResponse)
			}

			// Mengosongkan tanggal dummy jika ruangan tersedia
			if startTime == "1945-08-17" || *currentRoom.IsAvailable {
				startTime = ""
				endTime = ""
			}

			roomResponse = response.RoomResponse{
				Id:            currentRoom.Id,
				Name:          currentRoom.Name,
				RoomType:      currentRoom.RoomTypeId,
				Capacity:      currentRoom.Capacity,
				IsAvailable:   *currentRoom.IsAvailable,
				IsReserveable: *currentRoom.IsReserveable,
				StartTime:     startTime,
				EndTIme:       endTime,
			}
		}

		// Tambahkan fasilitas ke daftar jika fasilitas valid (bukan dummy "null")
		if currentFacility.Id != "null" {
			roomResponse.Facilities = append(roomResponse.Facilities, currentFacility)
		}
	}

	// Menambahkan ruangan terakhir ke list
	if roomResponse.Id != "" {
		responseData = append(responseData, roomResponse)
	}

	return responseData, nil
}

// NewRoomRepository menginisialisasi instansi konkrit RoomRepository baru.
func NewRoomRepository(db *sql.DB) RoomRepository {
	return &roomRepository{
		db: db,
	}
}
