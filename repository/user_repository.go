package repository

import (
	"database/sql"
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/model/dto/request"
	"go-roomify/utils"
	"go-roomify/utils/query"
)

type UserProfileRepository interface {
	GetList(page, size int) ([]model.UserProfile, dto.Paging, error)
	GetById(id string) (model.UserProfile, error)
	GetByUsername(username string) (model.UserProfile, error)
	Create(user request.UserProfileRequest) (request.UserProfileRequest, error)
	Update(user request.UserProfileRequest) (request.UserProfileRequest, error)
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
	qselect.Join("mst_division AS d", "up.division_id=d.id")
	qselect.Join("users AS u", "up.user_id=u.id")
	qselect.Join("mst_role AS r", "up.role_id=r.id")
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
	qcount := query.QSelect{DB: u.db}
	qcount.Table("mst_user_profile")
	qcount.Column("COUNT(id)")

	err = qcount.RunRow().Scan(&totalRows)
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
	qselect.Join("mst_division AS d", "up.division_id=d.id")
	qselect.Join("users AS u", "up.user_id=u.id")
	qselect.Join("mst_role AS r", "up.role_id=r.id")
	qselect.Where("up.id", "=", id)

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
	qselect.Join("mst_division AS d", "up.division_id=d.id")
	qselect.Join("users AS u", "up.user_id=u.id")
	qselect.Join("mst_role AS r", "up.role_id=r.id")
	qselect.Where("u.username", "=", username)

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

	_, err := qinsert.Run()
	if err != nil {
		return request.UserProfileRequest{}, err
	}

	return user, nil
}

func (u *userProfileRepository) Update(user request.UserProfileRequest) (request.UserProfileRequest, error) {
	qupdate := query.QUpdate{DB: u.db}

	qupdate.Table("mst_user_profile")
	qupdate.Set("full_name", user.FullName)
	qupdate.Set("division_id", user.DivisionId)
	qupdate.Set("address", user.Address)
	qupdate.Set("phone_number", user.PhoneNumber)
	qupdate.Set("role_id", user.RoleId)
	qupdate.Where("id", "=", user.Id)
	_, err := qupdate.Run()
	if err != nil {
		return request.UserProfileRequest{}, err
	}

	return user, nil
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
