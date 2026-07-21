// Package repository menangani akses dan manipulasi data langsung ke database PostgreSQL.
package repository

// Import package database/sql, model, dto, request, utils, dan query builder.
import (
	"database/sql"                 // Interface koneksi database SQL
	"go-roomify/model"             // Struct model data UserProfile
	"go-roomify/model/dto"         // DTO paginasi
	"go-roomify/model/dto/request" // DTO request profil pengguna
	"go-roomify/utils"             // Helper paginasi
	"go-roomify/utils/query"       // Builder query SQL dinamis
)

// UserProfileRepository merupakan kontrak interface untuk mengelola profil pengguna (nama, divisi, telepon, dll).
type UserProfileRepository interface {
	GetList(page, size int) ([]model.UserProfile, dto.Paging, error)           // Mengambil daftar profil pengguna berhalaman
	GetById(id string) (model.UserProfile, error)                             // Mengambil profil pengguna berdasarkan ID profil
	GetByUsername(username string) (model.UserProfile, error)                 // Mengambil profil pengguna berdasarkan username
	Create(user request.UserProfileRequest) (request.UserProfileRequest, error) // Menyimpan data profil pengguna baru
	Update(user request.UserProfileRequest) (request.UserProfileRequest, error) // Memperbarui data profil pengguna
	Delete(id string) error                                                   // Menghapus data profil pengguna berdasarkan ID
}

// userProfileRepository merupakan struktur konkrit pengelola data profil pengguna.
type userProfileRepository struct {
	db *sql.DB // Pointer koneksi database PostgreSQL
}

// GetList mengambil daftar profil pengguna lengkap dengan relasi Divisi, User Credential, dan Role secara paginasi.
func (u *userProfileRepository) GetList(page, size int) ([]model.UserProfile, dto.Paging, error) {
	// Menghitung offset baris data
	skip := (page - 1) * size

	// Inisialisasi query SELECT
	qselect := query.QSelect{DB: u.db}

	qselect.Table("mst_user_profile AS up")
	qselect.Column(
		"up.id",
		"up.full_name",
		"d.id",
		"d.name",
		"up.address",
		"up.phone_number",
		"u.id",
		"u.username",
		"u.password",
		"u.token",
		"r.id",
		"r.position",
	)
	// JOIN tabel relasi
	qselect.Join("mst_division AS d", "up.division_id=d.id")
	qselect.Join("users AS u", "up.user_id=u.id")
	qselect.Join("mst_role AS r", "up.role_id=r.id")
	qselect.Limit(size)
	qselect.Offset(skip)

	// Eksekusi query
	rows, err := qselect.Run()
	if err != nil {
		return nil, dto.Paging{}, err
	}

	// Slice penampung daftar profil pengguna
	var users []model.UserProfile
	for rows.Next() {
		var user model.UserProfile
		// Scan data ke nested struct (Division, User, Role)
		err = rows.Scan(
			&user.Id,
			&user.FullName,
			&user.Division.Id,
			&user.Division.Name,
			&user.Address,
			&user.PhoneNumber,
			&user.User.Id,
			&user.User.Username,
			&user.User.Password,
			&user.User.Token,
			&user.Role.Id,
			&user.Role.Position,
		)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		users = append(users, user)
	}

	// Menhitung total record profil
	var totalRows int
	qcount := query.QSelect{DB: u.db}
	qcount.Table("mst_user_profile")
	qcount.Column("COUNT(id)")

	err = qcount.RunRow().Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	// Buat DTO paginasi
	resultPagingDto := utils.Paginate(page, size, totalRows)
	return users, resultPagingDto, nil
}

// GetById mengambil detail profil pengguna tunggal berdasarkan ID.
func (u *userProfileRepository) GetById(id string) (model.UserProfile, error) {
	var user model.UserProfile
	qselect := query.QSelect{DB: u.db}

	qselect.Table("mst_user_profile AS up")
	qselect.Column(
		"up.id",
		"up.full_name",
		"d.id",
		"d.name",
		"up.address",
		"up.phone_number",
		"u.id",
		"u.username",
		"u.password",
		"u.token",
		"r.id",
		"r.position",
	)
	qselect.Join("mst_division AS d", "up.division_id=d.id")
	qselect.Join("users AS u", "up.user_id=u.id")
	qselect.Join("mst_role AS r", "up.role_id=r.id")
	qselect.Where("up.id", "=", id)

	// Scanning hasil baris tunggal
	err := qselect.RunRow().Scan(
		&user.Id,
		&user.FullName,
		&user.Division.Id,
		&user.Division.Name,
		&user.Address,
		&user.PhoneNumber,
		&user.User.Id,
		&user.User.Username,
		&user.User.Password,
		&user.User.Token,
		&user.Role.Id,
		&user.Role.Position,
	)
	if err != nil {
		return model.UserProfile{}, err
	}

	return user, nil
}

// GetByUsername mengambil profil pengguna berdasarkan username login.
func (u *userProfileRepository) GetByUsername(username string) (model.UserProfile, error) {
	var user model.UserProfile
	qselect := query.QSelect{DB: u.db}

	qselect.Table("mst_user_profile AS up")
	qselect.Column(
		"up.id",
		"up.full_name",
		"d.id",
		"d.name",
		"up.address",
		"up.phone_number",
		"u.id",
		"u.username",
		"u.password",
		"u.token",
		"r.id",
		"r.position",
	)
	qselect.Join("mst_division AS d", "up.division_id=d.id")
	qselect.Join("users AS u", "up.user_id=u.id")
	qselect.Join("mst_role AS r", "up.role_id=r.id")
	qselect.Where("u.username", "=", username)

	// Scanning hasil baris tunggal
	err := qselect.RunRow().Scan(
		&user.Id,
		&user.FullName,
		&user.Division.Id,
		&user.Division.Name,
		&user.Address,
		&user.PhoneNumber,
		&user.User.Id,
		&user.User.Username,
		&user.User.Password,
		&user.User.Token,
		&user.Role.Id,
		&user.Role.Position,
	)
	if err != nil {
		return model.UserProfile{}, err
	}

	return user, nil
}

// Create menyisipkan data profil pengguna baru ke tabel mst_user_profile.
func (u *userProfileRepository) Create(user request.UserProfileRequest) (request.UserProfileRequest, error) {
	qinsert := query.QInsert{DB: u.db}

	qinsert.Table("mst_user_profile")
	qinsert.Column(
		"id",
		"full_name",
		"division_id",
		"address",
		"phone_number",
		"user_id",
		"role_id",
	)
	qinsert.Values(
		user.Id,
		user.FullName,
		user.DivisionId,
		user.Address,
		user.PhoneNumber,
		user.UserId,
		user.RoleId,
	)

	// Eksekusi insert
	_, err := qinsert.Run()
	if err != nil {
		return request.UserProfileRequest{}, err
	}

	return user, nil
}

// Update memperbarui data profil pengguna (Nama, Divisi, Alamat, No HP, Role).
func (u *userProfileRepository) Update(user request.UserProfileRequest) (request.UserProfileRequest, error) {
	qupdate := query.QUpdate{DB: u.db}

	qupdate.Table("mst_user_profile")
	qupdate.Set("full_name", user.FullName)
	qupdate.Set("division_id", user.DivisionId)
	qupdate.Set("address", user.Address)
	qupdate.Set("phone_number", user.PhoneNumber)
	qupdate.Set("role_id", user.RoleId)
	qupdate.Where("id", "=", user.Id)

	// Eksekusi update
	_, err := qupdate.Run()
	if err != nil {
		return request.UserProfileRequest{}, err
	}

	return user, nil
}

// Delete menghapus record profil pengguna dari tabel mst_user_profile.
func (u *userProfileRepository) Delete(id string) error {
	qdelete := query.QDelete{DB: u.db}

	qdelete.Table("mst_user_profile")
	qdelete.Where("id", "=", id)

	// Eksekusi penghapusan
	_, err := qdelete.Run()
	if err != nil {
		return err
	}

	return nil
}

// NewUserProfileRepository menginisialisasi instansi konkrit UserProfileRepository baru.
func NewUserProfileRepository(db *sql.DB) UserProfileRepository {
	return &userProfileRepository{
		db: db,
	}
}
