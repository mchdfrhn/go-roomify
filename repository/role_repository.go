// Package repository menangani akses dan manipulasi data langsung ke database PostgreSQL.
package repository

// Import package database/sql, model, dto, utils, dan query builder.
import (
	"database/sql"           // Interface koneksi database SQL
	"go-roomify/model"       // Struct model data Role
	"go-roomify/model/dto"   // DTO paginasi
	"go-roomify/utils"       // Helper paginasi
	"go-roomify/utils/query" // Builder query SQL dinamis
)

// RoleRepository merupakan kontrak interface untuk mengelola data master Peran/Jabatan (Role).
type RoleRepository interface {
	InsertRole(newRole model.Role) error                                 // Menyimpan data role baru
	UpdateRole(newRole model.Role) error                                 // Memperbarui posisi role
	GetRoleById(id string) (model.Role, error)                           // Mengambil data role berdasarkan ID
	DeleteRoleById(id string) error                                      // Menghapus data role berdasarkan ID
	GetListPaging(page int, size int) ([]model.Role, dto.Paging, error) // Mengambil daftar role berhalaman
}

// roleRepository merupakan struktur konkrit pengelola data master role.
type roleRepository struct {
	db *sql.DB // Pointer koneksi database PostgreSQL
}

// GetListPaging mengambil daftar role berhalaman menggunakan Limit dan Offset.
func (r *roleRepository) GetListPaging(page int, size int) ([]model.Role, dto.Paging, error) {
	// Menhitung jumlah baris yang di-skip
	skip := (page - 1) * size
	// Inisialisasi builder SELECT
	qSelect := &query.QSelect{DB: r.db}
	qSelect.Table("mst_role").
		Column("id", "position").
		Limit(size).
		Offset(skip)

	// Eksekusi query SELECT
	rows, err := qSelect.Run()
	if err != nil {
		return nil, dto.Paging{}, err
	}
	// Memastikan resource rows ditutup setelah fungsi selesai
	defer rows.Close()

	// Slice penampung daftar role
	var roles []model.Role
	// Iterasi hasil query
	for rows.Next() {
		var role model.Role
		// Scan nilai kolom id dan position
		err = rows.Scan(
			&role.Id,
			&role.Position,
		)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		roles = append(roles, role)
	}

	// Menghitung total data role di tabel
	var totalRows int
	err = r.db.QueryRow("SELECT COUNT(*) FROM mst_role").Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	// Mengolah DTO metadata paginasi
	resultPagingDto := utils.Paginate(page, size, totalRows)
	return roles, resultPagingDto, nil
}

// UpdateRole memperbarui nama posisi role berdasarkan ID.
func (r *roleRepository) UpdateRole(newRole model.Role) error {
	// Inisialisasi builder UPDATE
	qUpdate := &query.QUpdate{DB: r.db}
	qUpdate.Table("mst_role").
		Set("position", newRole.Position).
		Where("id", "=", newRole.Id)

	// Eksekusi update
	_, err := qUpdate.Run()
	if err != nil {
		return err
	}
	return nil
}

// GetRoleById mengambil satu record role berdasarkan ID (UUID).
func (r *roleRepository) GetRoleById(id string) (model.Role, error) {
	// Inisialisasi builder SELECT
	qSelect := &query.QSelect{DB: r.db}
	qSelect.Table("mst_role").
		Column("id", "position").
		Where("id", "=", id)

	// Eksekusi query baris tunggal
	row := qSelect.RunRow()
	var role model.Role
	// Scan data ke struct role
	err := row.Scan(
		&role.Id,
		&role.Position,
	)

	if err != nil {
		return model.Role{}, err
	}
	return role, nil
}

// DeleteRoleById menghapus record role berdasarkan ID.
func (r *roleRepository) DeleteRoleById(id string) error {
	// Inisialisasi builder DELETE
	qDelete := &query.QDelete{DB: r.db}
	qDelete.Table("mst_role").Where("id", "=", id)

	// Eksekusi penghapusan
	_, err := qDelete.Run()
	if err != nil {
		return err
	}
	return nil
}

// newRoleRepository (Private constructor helper).
func newRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

// InsertRole menyisipkan data role baru ke tabel mst_role.
func (r *roleRepository) InsertRole(newRole model.Role) error {
	// Inisialisasi builder INSERT
	qInsert := query.QInsert{DB: r.db}
	qInsert.Table("mst_role").
		Column("id", "position").
		Values(newRole.Id, newRole.Position)

	// Eksekusi penyisipan data
	_, err := qInsert.Run()

	if err != nil {
		return err
	}
	return nil
}

// NewRoleRepository menginisialisasi instansi konkrit RoleRepository baru.
func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}
