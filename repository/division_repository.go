// Package repository menangani akses dan manipulasi data langsung ke database PostgreSQL.
package repository

// Import package database/sql, model, dto, utils, dan query builder.
import (
	"database/sql"           // Interface koneksi database SQL
	"go-roomify/model"       // Entitas model divisi
	"go-roomify/model/dto"   // DTO paginasi
	"go-roomify/utils"       // Helper paginasi
	"go-roomify/utils/query" // Builder query SQL dinamis
)

// DivisionRepository merupakan kontrak interface operasi CRUD pada data divisi.
type DivisionRepository interface {
	GetAllDivisi(page, size int) ([]model.Division, dto.Paging, error) // Mengambil seluruh divisi dengan paginasi
	GetDivisiById(id string) (model.Division, error)                  // Mengambil detail divisi berdasarkan ID
	CreateDivisi(division model.Division) error                        // Menyimpan divisi baru
	UpdateDivisi(division model.Division) error                        // Memperbarui data divisi
	DeleteDivisi(id string) error                                      // Menghapus data divisi berdasarkan ID
}

// divisionRepository merupakan struktur konkrit pengelola repository divisi.
type divisionRepository struct {
	db *sql.DB // Pointer koneksi database PostgreSQL
}

// GetAllDivisi mengambil daftar divisi berhalaman beserta informasi metadata paginasi.
func (r *divisionRepository) GetAllDivisi(page, size int) ([]model.Division, dto.Paging, error) {
	// Menhitung offset baris data yang dilewati
	skip := (page - 1) * size
	// Inisialisasi builder SELECT
	qselect := query.QSelect{DB: r.db}

	// Menentukan tabel target mst_division
	qselect.Table("mst_division")
	// Menentukan kolom yang dipilih: id, name
	qselect.Column(
		"id",
		"name")
	// Mengatur batas jumlah baris
	qselect.Limit(size)
	// Mengatur offset baris
	qselect.Offset(skip)
	// Eksekusi query SELECT
	rows, err := qselect.Run()
	// Memeriksa jika eksekusi gagal
	if err != nil {
		return nil, dto.Paging{}, err
	}
	// Memastikan cursor rows ditutup setelah selesai
	defer rows.Close()

	// Inisialisasi slice penampung hasil divisi
	var divisions []model.Division
	// Iterasi setiap baris hasil query
	for rows.Next() {
		var division model.Division
		// Memindahkan nilai kolom ke struct division
		if err := rows.Scan(
			&division.Id,
			&division.Name); err != nil {
			return nil, dto.Paging{}, err
		}
		// Menambahkan divisi ke slice
		divisions = append(divisions, division)
	}
	// Inisialisasi variabel hitung total data
	var totalRows int
	// Builder SELECT untuk COUNT data
	qcount := query.QSelect{DB: r.db}
	qcount.Table("mst_division")
	qcount.Column("COUNT(id)")

	// Eksekusi penanganan jumlah total data divisi
	err = qcount.RunRow().Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	// Mengalkulasi metadata paginasi
	resultPagingDto := utils.Paginate(page, size, totalRows)
	// Mengembalikan daftar divisi, DTO paginasi, dan nil error
	return divisions, resultPagingDto, nil
}

// GetDivisiById mengambil satu baris divisi berdasarkan ID unik.
func (r *divisionRepository) GetDivisiById(id string) (model.Division, error) {
	// Inisialisasi builder SELECT
	qselect := query.QSelect{DB: r.db}

	// Menyusun query SELECT mst_division WHERE id = id
	qselect.Table("mst_division")
	qselect.Column("id", "name")
	qselect.Where("id", "=", id)

	// Eksekusi single row query
	row := qselect.RunRow()
	var division model.Division
	// Memindahkan data kolom ke atribut struct division
	err := row.Scan(&division.Id, &division.Name)
	if err != nil {
		return division, err
	}

	// Mengembalikan objek divisi dan nil error
	return division, nil
}

// CreateDivisi menyisipkan record data divisi baru ke tabel mst_division.
func (r *divisionRepository) CreateDivisi(division model.Division) error {
	// Inisialisasi builder INSERT
	qinsert := query.QInsert{DB: r.db}
	// Menyusun query INSERT INTO mst_division (id, name) VALUES (...)
	qinsert.Table("mst_division").Column("id", "name").Values(division.Id, division.Name)
	// Eksekusi perintah INSERT
	_, err := qinsert.Run()
	// Mengembalikan error jika ada
	return err
}

// UpdateDivisi memperbarui nama divisi berdasarkan ID pada tabel mst_division.
func (r *divisionRepository) UpdateDivisi(division model.Division) error {
	// Inisialisasi builder UPDATE
	qUpdate := query.QUpdate{DB: r.db}
	// Menyusun query UPDATE mst_division SET name = ... WHERE id = ...
	qUpdate.Table("mst_division").
		Set("name", division.Name).
		Where("id", "=", division.Id)
	// Eksekusi pembaruan data
	_, err := qUpdate.Run()
	// Mengembalikan error jika ada
	return err
}

// DeleteDivisi menghapus data divisi berdasarkan ID dari tabel mst_division.
func (r *divisionRepository) DeleteDivisi(id string) error {
	// Inisialisasi builder DELETE
	qdelete := query.QDelete{DB: r.db}
	// Menyusun query DELETE FROM mst_division WHERE id = ...
	qdelete.Table("mst_division").Where("id", "=", id)
	// Eksekusi perintah penghapusan
	_, err := qdelete.Run()
	// Mengembalikan error jika ada
	return err
}

// NewDivisionRepository menginstansiasi instance konkrit DivisionRepository baru.
func NewDivisionRepository(db *sql.DB) DivisionRepository {
	// Mengembalikan pointer divisi repository
	return &divisionRepository{db: db}
}
