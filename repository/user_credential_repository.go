package repository

import (
	"database/sql"
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/model/dto/request"
	"go-roomify/utils"
	"go-roomify/utils/query"
)

type UserCredentialRepository interface {
	GetByUsername(username string) (model.UserCredentialJwt, error)
	GetById(id string) (model.UserCredential, error)
	GetList(page int, size int) ([]model.UserCredential, dto.Paging, error)
	AddNew(new_user model.UserCredential) (model.UserCredential, error)
	UpdatePassword(new_user request.UserUpdatePasswordRequest) error
	Delete(id string) error
}

type userCredentialRepository struct {
	db *sql.DB
}

func (self *userCredentialRepository) GetByUsername(username string) (model.UserCredentialJwt, error) {
	qselect := query.QSelect{DB: self.db}

	qselect.Table("mst_user_profile")
	qselect.Column(
		"users.id",
		"users.username",
		"users.password",
		"mst_role.position",
		"users.token")
	qselect.Join("mst_role", "mst_user_profile.role_id = mst_role.id")
	qselect.Join("users", "mst_user_profile.user_id = users.id")
	qselect.Where("users.username", "=", username)
	qselect.Limit(1)

	rows, err := qselect.Run()

	if err != nil {
		return model.UserCredentialJwt{}, err
	}

	var r_user_cr model.UserCredentialJwt

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

	rows.Close()

	return r_user_cr, nil
}

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

func (self *userCredentialRepository) GetList(page int, size int) ([]model.UserCredential, dto.Paging, error) {
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

	rows, err := qselect.Run()
	if err != nil {
		return nil, dto.Paging{}, err
	}

	var rows_user []model.UserCredential

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

	var total_rows int
	qcount := query.QSelect{DB: self.db}
	qcount.Table("users")
	qcount.Column("COUNT(id)")

	err = qcount.RunRow().Scan(&total_rows)

	if err != nil {
		return nil, dto.Paging{}, err
	}

	resultPagingDto := utils.Paginate(page, size, total_rows)

	return rows_user, resultPagingDto, nil
}

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

	_, err := qinsert.Run()

	if err != nil {
		return model.UserCredential{}, err
	}

	return new_user, nil
}

func (self *userCredentialRepository) UpdatePassword(new_user request.UserUpdatePasswordRequest) error {
	qupdate := query.QUpdate{DB: self.db}

	qupdate.Table("users")
	qupdate.Set("password", new_user.NewPassword)
	qupdate.Where("id", "=", new_user.Id)

	_, err := qupdate.Run()

	return err
}

func (self *userCredentialRepository) Delete(id string) error {
	qdelete := query.QDelete{DB: self.db}

	qdelete.Table("users")
	qdelete.Where("id", "=", id)

	_, err := qdelete.Run()

	return err
}

func NewUserCredentialRepository(db *sql.DB) UserCredentialRepository {
	return &userCredentialRepository{
		db: db,
	}
}
