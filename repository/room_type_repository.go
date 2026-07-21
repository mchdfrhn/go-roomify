// Package repository menangani akses dan manipulasi data langsung ke database PostgreSQL.
package repository

// Import package database/sql, model data, dan query builder.
import (
	"database/sql"           // Interface koneksi database SQL
	"go-roomify/model"       // Struct model data RoomType
	"go-roomify/utils/query" // Builder query SQL dinamis
)

// RoomTypeRepository merupakan kontrak interface untuk mengelola data master Kategori/Tipe Ruangan.
type RoomTypeRepository interface {
	CreateRoomType(roomModelType model.RoomType) error                       // Menyimpan kategori ruangan baru
	GetAllRoomType() ([]model.RoomType, error)                               // Mengambil semua daftar tipe ruangan
	GetRoomTypeByIdOrName(idOrNameRoomType string) (model.RoomType, error)   // Mencari tipe ruangan berdasarkan ID atau nama
	UpdateRoomTypeById(updateRoomType model.RoomType) (model.RoomType, error) // Memperbarui nama tipe ruangan
	DeleteRoomTypeById(roomId string) error                                  // Menghapus tipe ruangan berdasarkan ID
}

// roomTypeRepository merupakan struktur konkrit pengelola data master tipe ruangan.
type roomTypeRepository struct {
	db *sql.DB // Pointer koneksi database PostgreSQL
}

// CreateRoomType menyisipkan record tipe ruangan baru ke tabel room_type.
func (rr *roomTypeRepository) CreateRoomType(roomModelType model.RoomType) error {
	// Inisialisasi builder INSERT
	query := query.QInsert{DB: rr.db}

	// Eksekusi penyisipan nama dan id tipe ruangan
	_, err := query.Table(
		"room_type",
	).Column(
		"id", "name",
	).Values(
		roomModelType.Id, roomModelType.Name,
	).Run()

	if err != nil {
		return err
	}

	return nil
}

// GetAllRoomType mengambil seluruh data kategori ruangan tanpa paginasi.
func (rr *roomTypeRepository) GetAllRoomType() ([]model.RoomType, error) {
	// Inisialisasi builder SELECT
	queryAllRoom := query.QSelect{DB: rr.db}

	// Menjalankan SELECT id, name FROM room_type
	rows, err := queryAllRoom.Table(
		"room_type",
	).Column(
		"id", "name",
	).Run()

	if err != nil {
		return nil, err
	}

	// Slice penampung daftar tipe ruangan
	var allRoomType []model.RoomType
	for rows.Next() {
		var roomType model.RoomType

		// Scan data kolom ke struct
		err := rows.Scan(
			&roomType.Id, &roomType.Name,
		)

		if err != nil {
			return nil, err
		}

		allRoomType = append(allRoomType, roomType)
	}

	return allRoomType, nil
}

// GetRoomTypeByIdOrName mencari tipe ruangan tunggal berdasarkan kecocokan ID atau Nama.
func (rr *roomTypeRepository) GetRoomTypeByIdOrName(roomTypeIdOrName string) (model.RoomType, error) {
	// Inisialisasi builder SELECT
	query := query.QSelect{DB: rr.db}
	var roomType model.RoomType

	// Query dengan klausul WHERE id = x OR name = x
	err := query.Table(
		"room_type",
	).Column(
		"id",
		"name",
	).Where(
		"id", "=", roomTypeIdOrName,
	).OrWhere(
		"name", "=", roomTypeIdOrName,
	).RunRow().Scan(
		&roomType.Id, &roomType.Name,
	)

	if err != nil {
		// Mengabaikan error jika data tidak ditemukan (ErrNoRows)
		if err == sql.ErrNoRows {
			return model.RoomType{}, nil
		}
		return model.RoomType{}, err
	}

	return roomType, nil
}

// UpdateRoomTypeById memperbarui nama tipe ruangan berdasarkan ID.
func (rr *roomTypeRepository) UpdateRoomTypeById(updateRoomType model.RoomType) (model.RoomType, error) {
	// Inisialisasi builder UPDATE
	query := query.QUpdate{DB: rr.db}

	// Eksekusi update nama
	_, err := query.Table(
		"room_type",
	).Set(
		"name", updateRoomType.Name,
	).Where(
		"id", "=", updateRoomType.Id,
	).Run()

	if err != nil {
		return model.RoomType{}, err
	}

	return updateRoomType, nil
}

// DeleteRoomTypeById menghapus record tipe ruangan dari database.
func (rr *roomTypeRepository) DeleteRoomTypeById(roomId string) error {
	// Inisialisasi builder DELETE
	query := query.QDelete{DB: rr.db}

	// Eksekusi penghapusan
	_, err := query.Table(
		"room_type",
	).Where(
		"id", "=", roomId,
	).Run()

	if err != nil {
		return err
	}

	return nil
}

// NewRoomTypeRepository menginisialisasi instansi konkrit RoomTypeRepository baru.
func NewRoomTypeRepository(db *sql.DB) RoomTypeRepository {
	return &roomTypeRepository{
		db: db,
	}
}
