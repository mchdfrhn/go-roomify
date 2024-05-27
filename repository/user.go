package repository

import (
	"database/sql"
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/utils"
	"go-roomify/utils/query"
)

type UserRepository interface {
	GetList(page, size int) ([]model.User, dto.Paging, error)
	GetById(id string) (model.User, error)
	GetByUsername(username string) (model.User, error)
	Update(payload model.User) (model.User, error)
	Delete(id string) error
}

type userRepository struct {
	db *sql.DB
}

func (u *userRepository) GetList(page, size int) ([]model.User, dto.Paging, error) {
	skip := (page - 1) * size

	qselect := query.QSelect{DB: u.db}

	qselect.Table("users")
	qselect.Column(
		"id",
		"username",
		"password",
		"token",
	)
	qselect.Limit(size)
	qselect.Offset(skip)

	rows, err := qselect.Run()
	if err != nil {
		return nil, dto.Paging{}, err
	}

	var users []model.User
	for rows.Next() {
		var user model.User
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.Password,
			&user.Token,
		)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		users = append(users, user)
	}
	var totalRows int
	qselect.Table("users")
	qselect.Column("COUNT (*)")

	err = qselect.RunRow().Scan(totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	resultPagingDto := utils.Paginate(page, size, totalRows)
	return users, resultPagingDto, nil
}

func (u *userRepository) GetById(id string) (model.User, error) {
	var user model.User
	qselect := query.QSelect{DB: u.db}

	qselect.Table("users")
	qselect.Column(
		"id",
		"username",
		"password",
		"token",
	)
	qselect.Where("id", "=", id)

	err := qselect.RunRow().Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Token,
	)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *userRepository) GetByUsername(username string) (model.User, error) {
	var user model.User
	qselect := query.QSelect{DB: u.db}

	qselect.Table("users")
	qselect.Column(
		"id",
		"username",
		"password",
		"token",
	)
	qselect.Where("username", "=", username)

	err := qselect.RunRow().Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Token,
	)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *userRepository) Update(payload model.User) (model.User, error) {
	qupdate := query.QUpdate{DB: u.db}

	qupdate.Table("users")
	qupdate.Set("id", payload.Id)
	qupdate.Set("username", payload.Username)
	qupdate.Set("password", payload.Password)
	qupdate.Set("token", payload.Token)

	_, err := qupdate.Run()
	if err != nil {
		return model.User{}, err
	}

	return payload, nil
}

func (u *userRepository) Delete(id string) error {
	qdelete := query.QDelete{DB: u.db}

	qdelete.Table("users")
	qdelete.Where("id", "=", id)

	_, err := qdelete.Run()
	if err != nil {
		return err
	}

	return nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
