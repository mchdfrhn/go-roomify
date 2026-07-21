// Package repository menangani akses dan manipulasi data langsung ke database PostgreSQL.
package repository

// Import package database/sql, model, dto, utils, query builder, dan strconv.
import (
	"database/sql"           // Interface koneksi database SQL
	"go-roomify/model"       // Model data fasilitas
	"go-roomify/model/dto"   // DTO paginasi
	"go-roomify/utils"       // Utility paginasi
	"go-roomify/utils/query" // Builder query SQL dinamis
	"strconv"                // Konversi tipe boolean ke string
)

// FacilityRepository merupakan kontrak interface operasi data fasilitas ruangan.
type FacilityRepository interface {
	GetFacilityById(id string) (model.Facility, error)                         // Mengambil fasilitas berdasarkan ID
	GetPagingFacility(page int, size int) ([]model.Facility, dto.Paging, error) // Mengambil daftar fasilitas berhalaman
	InsertFacility(newFacility model.Facility) error                           // Menyimpan data fasilitas baru
	UpdateFacility(newFacility model.Facility) error                           // Memperbarui data fasilitas
	DeleteFacility(id string) error                                            // Menghapus fasilitas berdasarkan ID
}

// facilityRepository merupakan struktur konkrit implementasi FacilityRepository.
type facilityRepository struct {
	db *sql.DB // Pointer koneksi database PostgreSQL
}

// NewFacilityRepository membuat instansi facilityRepository baru.
func NewFacilityRepository(db *sql.DB) FacilityRepository {
	return &facilityRepository{
		db: db,
	}
}

// GetFacilityById mencari satu record fasilitas berdasarkan ID (UUID).
func (f *facilityRepository) GetFacilityById(id string) (model.Facility, error) {
	var facility model.Facility
	// Inisialisasi builder SELECT
	qselect := query.QSelect{DB: f.db}
	qselect.Table("mst_facility")
	// Menentukan kolom pencarian
	qselect.Column(
		"id",
		"name",
		"is_available",
		"is_reserveable",
		"room_id",
	)
	// Kondisi WHERE id = id
	qselect.Where("id", "=", id)
	// Eksekusi baris tunggal
	rows := qselect.RunRow()
	// Scan nilai kolom ke struct facility
	err := rows.Scan(
		&facility.Id,
		&facility.Name,
		&facility.IsAvailable,
		&facility.IsReserveable,
		&facility.RoomId,
	)

	// Jika terjadi error (data tidak ditemukan)
	if err != nil {
		return model.Facility{}, err
	}

	// Mengembalikan data fasilitas
	return facility, nil
}

// GetPagingFacility mengambil daftar fasilitas secara berhalaman (paginated).
func (f *facilityRepository) GetPagingFacility(page int, size int) ([]model.Facility, dto.Paging, error) {
	// Menghitung jumlah baris yang di-skip
	skip := (page - 1) * size
	// Inisialisasi builder query SELECT
	qselect := query.QSelect{DB: f.db}
	qselect.Table("mst_facility")
	qselect.Column(
		"id",
		"name",
		"is_available",
		"is_reserveable",
		"room_id",
	)
	qselect.Limit(size)
	qselect.Offset(skip)

	// Eksekusi query SELECT
	rows, err := qselect.Run()
	if err != nil {
		return nil, dto.Paging{}, err
	}

	// Slice penampung daftar fasilitas
	var facilities []model.Facility
	// Iterasi baris data
	for rows.Next() {
		var facility model.Facility
		// Memindahkan data kolom
		err = rows.Scan(
			&facility.Id,
			&facility.Name,
			&facility.IsAvailable,
			&facility.IsReserveable,
			&facility.RoomId,
		)
		if err != nil {
			return nil, dto.Paging{}, err
		}

		// Menambahkan fasilitas ke slice
		facilities = append(facilities, facility)
	}

	// Mengalkulasi jumlah seluruh record fasilitas
	var totalRows int
	qcount := query.QSelect{DB: f.db}
	qcount.Table("mst_facility")
	qcount.Column("COUNT(id)")

	err = qcount.RunRow().Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	// Mengolah DTO paginasi
	resultPagingDto := utils.Paginate(page, size, totalRows)
	return facilities, resultPagingDto, nil
}

// InsertFacility menyisipkan data fasilitas baru ke tabel mst_facility.
func (f *facilityRepository) InsertFacility(newFacility model.Facility) error {
	// Inisialisasi builder INSERT
	qinsert := query.QInsert{DB: f.db}
	qinsert.Table("mst_facility")
	qinsert.Column(
		"id",
		"name",
		"is_available",
		"is_reserveable",
		"room_id",
	)
	// Binding nilai atribut fasilitas
	qinsert.Values(
		newFacility.Id,
		newFacility.Name,
		strconv.FormatBool(newFacility.IsAvailable),
		strconv.FormatBool(newFacility.IsReserveable),
		newFacility.RoomId,
	)

	// Eksekusi query penyisipan data
	_, err := qinsert.Run()
	if err != nil {
		return err
	}

	return nil
}

// UpdateFacility memperbarui informasi fasilitas yang sudah ada.
func (f *facilityRepository) UpdateFacility(newFacility model.Facility) error {
	// Inisialisasi builder UPDATE
	qupdate := query.QUpdate{DB: f.db}
	qupdate.Table("mst_facility")
	// Menentukan set kolom yang diubah
	qupdate.Set("name", newFacility.Name)
	qupdate.Set("is_available", strconv.FormatBool(newFacility.IsAvailable))
	qupdate.Set("is_reserveable", strconv.FormatBool(newFacility.IsReserveable))
	qupdate.Set("room_id", newFacility.RoomId)
	// Kondisi kriteria update WHERE id = id
	qupdate.Where("id", "=", newFacility.Id)

	// Eksekusi perbaruan data
	_, err := qupdate.Run()
	if err != nil {
		return err
	}

	return nil
}

// DeleteFacility menghapus data fasilitas berdasarkan ID.
func (f *facilityRepository) DeleteFacility(id string) error {
	// Inisialisasi builder DELETE
	qdelete := query.QDelete{DB: f.db}
	qdelete.Table("mst_facility")
	qdelete.Where("id", "=", id)
	// Eksekusi penghapusan
	_, err := qdelete.Run()
	if err != nil {
		return err
	}

	return nil
}
