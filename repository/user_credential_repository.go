// Package repository menangani akses dan manipulasi data langsung ke database PostgreSQL.
package repository

// Import package database/sql, model, dto, request, utils, dan query builder.
import (
	"database/sql"                 // Interface koneksi database SQL
	"go-roomify/model"             // Struct model data UserCredential & UserCredentialJwt
	"go-roomify/model/dto"         // DTO paginasi
	"go-roomify/model/dto/request" // DTO request perbaruan password
	"go-roomify/utils"             // Helper paginasi
	"go-roomify/utils/query"       // Builder query SQL dinamis
)

// UserCredentialRepository merupakan kontrak interface untuk mengelola kredensial login akun pengguna.
type UserCredentialRepository interface {
	GetByUsername(username string) (model.UserCredentialJwt, error)          // Mengambil data kredensial + role JWT berdasarkan username
	GetById(id string) (model.UserCredential, error)                        // Mengambil kredensial berdasarkan ID
	GetList(page int, size int) ([]model.UserCredential, dto.Paging, error) // Mengambil daftar kredensial secara paginasi
	AddNew(new_user model.UserCredential) (model.UserCredential, error)     // Menyimpan kredensial baru
	UpdatePassword(new_user request.UserUpdatePasswordRequest) error        // Memperbarui hash password pengguna
	Delete(id string) error                                                 // Menghapus akun kredensial berdasarkan ID
}

// userCredentialRepository merupakan struktur konkrit pengelola kredensial akun.
type userCredentialRepository struct {
	db *sql.DB // Pointer koneksi database PostgreSQL
}

// GetByUsername mengambil data login pengguna beserta nama peran (role) untuk klaim JWT.
func (self *userCredentialRepository) GetByUsername(username string) (model.UserCredentialJwt, error) {
	// Inisialisasi builder SELECT
	qselect := query.QSelect{DB: self.db}

	qselect.Table("mst_user_profile")
	qselect.Column(
		"users.id",
		"users.username",
		"users.password",
		"mst_role.position",
		"users.token")
	// JOIN tabel mst_role dan users
	qselect.Join("mst_role", "mst_user_profile.role_id = mst_role.id")
	qselect.Join("users", "mst_user_profile.user_id = users.id")
	qselect.Where("users.username", "=", username)
	qselect.Limit(1)

	// Eksekusi query SELECT
	rows, err := qselect.Run()

	if err != nil {
		return model.UserCredentialJwt{}, err
	}

	var r_user_cr model.UserCredentialJwt

	// Iterasi hasil scanning
	for rows.Next() {
		err := rows.Scan(
			&r_user_cr.Id,
			&r_user_cr.Username,
			&r_user_cr.Password,
			&r_user_cr.Role,
			&r_user_cr.Token)

		if err != nil {
			return model.UserCredentialJwt{}, err
		}
	}

	// Menutup penampung rows
	rows.Close()

	return r_user_cr, nil
}

// GetById mengambil record kredensial tunggal dari tabel users berdasarkan ID.
func (self *userCredentialRepository) GetById(id string) (model.UserCredential, error) {
	qselect := query.QSelect{DB: self.db}

	qselect.Table("users")
	qselect.Column(
		"id",
		"username",
		"password",
		"token")
	qselect.Where("id", "=", id)

	var r_user_cr model.UserCredential

	// Eksekusi query baris tunggal
	err := qselect.RunRow().Scan(
		&r_user_cr.Id,
		&r_user_cr.Username,
		&r_user_cr.Password,
		&r_user_cr.Token)

	if err != nil {
		return model.UserCredential{}, err
	}

	return r_user_cr, nil
}

// GetList mengambil daftar akun pengguna berhalaman.
func (self *userCredentialRepository) GetList(page int, size int) ([]model.UserCredential, dto.Paging, error) {
	// Menghitung offset data
	skip := (page - 1) * size

	qselect := query.QSelect{DB: self.db}
	qselect.Table("users")
	qselect.Column(
		"id",
		"username",
		"password",
		"token")
	qselect.Limit(size)
	qselect.Offset(skip)

	// Eksekusi query
	rows, err := qselect.Run()
	if err != nil {
		return nil, dto.Paging{}, err
	}

	var rows_user []model.UserCredential

	// Iterasi scanning baris data pengguna
	for rows.Next() {
		var r_user_cr model.UserCredential

		err := rows.Scan(
			&r_user_cr.Id,
			&r_user_cr.Username,
			&r_user_cr.Password,
			&r_user_cr.Token)

		if err != nil {
			return nil, dto.Paging{}, err
		}

		rows_user = append(rows_user, r_user_cr)
	}

	// Hitung total jumlah record users
	var total_rows int
	qcount := query.QSelect{DB: self.db}
	qcount.Table("users")
	qcount.Column("COUNT(id)")

	err = qcount.RunRow().Scan(&total_rows)

	if err != nil {
		return nil, dto.Paging{}, err
	}

	// Buat DTO paginasi
	resultPagingDto := utils.Paginate(page, size, total_rows)

	return rows_user, resultPagingDto, nil
}

// AddNew menyisipkan akun kredensial baru (username & hashed password) ke tabel users.
func (self *userCredentialRepository) AddNew(new_user model.UserCredential) (model.UserCredential, error) {
	qinsert := query.QInsert{DB: self.db}

	qinsert.Table("users")
	qinsert.Column(
		"id",
		"username",
		"password",
		"token",
	)
	qinsert.Values(
		new_user.Id,
		new_user.Username,
		new_user.Password,
		new_user.Token,
	)

	// Eksekusi insert
	_, err := qinsert.Run()

	if err != nil {
		return model.UserCredential{}, err
	}

	return new_user, nil
}

// UpdatePassword memperbarui nilai hash password milik pengguna di tabel users.
func (self *userCredentialRepository) UpdatePassword(new_user request.UserUpdatePasswordRequest) error {
	qupdate := query.QUpdate{DB: self.db}

	qupdate.Table("users")
	qupdate.Set("password", new_user.NewPassword)
	qupdate.Where("id", "=", new_user.Id)

	// Eksekusi perbaruan password
	_, err := qupdate.Run()

	return err
}

// Delete menghapus record kredensial akun dari tabel users.
func (self *userCredentialRepository) Delete(id string) error {
	qdelete := query.QDelete{DB: self.db}

	qdelete.Table("users")
	qdelete.Where("id", "=", id)

	// Eksekusi penghapusan
	_, err := qdelete.Run()

	return err
}

// NewUserCredentialRepository menginisialisasi instansi konkrit UserCredentialRepository baru.
func NewUserCredentialRepository(db *sql.DB) UserCredentialRepository {
	return &userCredentialRepository{
		db: db,
	}
}
