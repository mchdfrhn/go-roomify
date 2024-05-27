package repository

import (
	"database/sql"
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/utils"
	"go-roomify/utils/query"
)

type UserProfileRepository interface {
	GetList(page, size int) ([]model.UserProfile, dto.Paging, error)
	GetById(id string) (model.UserProfile, error)
	GetByUsername(username string) (model.UserProfile, error)
	Create(payload model.UserProfile) (model.UserProfile, error)
	Update(payload model.UserProfile) (model.UserProfile, error)
	Delete(id string) error
}

type userProfileRepository struct {
	db *sql.DB
}

func (u *userProfileRepository) GetList(page, size int) ([]model.UserProfile, dto.Paging, error) {
	skip := (page - 1) * size

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
	qselect.Join("mst_division AS d", "ON up.division_id=d.id")
	qselect.Join("users AS u", "ON up.user_id=u.id")
	qselect.Join("mst_role AS r", "ON up.role_id=r.id")
	qselect.Limit(size)
	qselect.Offset(skip)

	rows, err := qselect.Run()
	if err != nil {
		return nil, dto.Paging{}, err
	}

	var users []model.UserProfile
	for rows.Next() {
		var user model.UserProfile
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

	var totalRows int
	qselect.Table("mst_user_profile")
	qselect.Column("COUNT (*)")

	err = qselect.RunRow().Scan(totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	resultPagingDto := utils.Paginate(page, size, totalRows)
	return users, resultPagingDto, nil
}

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
	qselect.Join("mst_division AS d", "ON up.division_id=d.id")
	qselect.Join("users AS u", "ON up.user_id=u.id")
	qselect.Join("mst_role AS r", "ON up.role_id=r.id")
	qselect.Where("id", "=", id)

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
	qselect.Join("mst_division AS d", "ON up.division_id=d.id")
	qselect.Join("users AS u", "ON up.user_id=u.id")
	qselect.Join("mst_role AS r", "ON up.role_id=r.id")
	qselect.Where("username", "=", username)

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

func (u *userProfileRepository) Create(payload model.UserProfile) (model.UserProfile, error) {
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
		payload.Id,
		payload.FullName,
		payload.Division.Id,
		payload.Address,
		payload.PhoneNumber,
		payload.User.Id,
		payload.Role.Id,
	)

	_, err := qinsert.Run()
	if err != nil {
		return model.UserProfile{}, err
	}

	return payload, nil
}

func (u *userProfileRepository) Update(payload model.UserProfile) (model.UserProfile, error) {
	qupdate := query.QUpdate{DB: u.db}

	qupdate.Table("mst_user_profile")
	qupdate.Set("full_name", payload.FullName)
	qupdate.Set("division_id", payload.Division.Id)
	qupdate.Set("address", payload.Address)
	qupdate.Set("phone_number", payload.PhoneNumber)
	qupdate.Set("role_id", payload.Role.Id)

	_, err := qupdate.Run()
	if err != nil {
		return model.UserProfile{}, err
	}

	return payload, nil
}

func (u *userProfileRepository) Delete(id string) error {
	qdelete := query.QDelete{DB: u.db}

	qdelete.Table("mst_user_profile")
	qdelete.Where("id", "=", id)

	_, err := qdelete.Run()
	if err != nil {
		return err
	}

	return nil
}

func NewUserProfileRepository(db *sql.DB) UserProfileRepository {
	return &userProfileRepository{
		db: db,
	}
}
