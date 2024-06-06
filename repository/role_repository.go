package repository

import (
	"database/sql"
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/utils"
	"go-roomify/utils/query"
)

type RoleRepository interface {
	InsertRole(newRole model.Role) error
	UpdateRole(newRole model.Role) error
	GetRoleById(id string) (model.Role, error)
	DeleteRoleById(id string) error
	GetListPaging(page int, size int) ([]model.Role, dto.Paging, error)
}

type roleRepository struct {
	db *sql.DB
}

func (r *roleRepository) GetListPaging(page int, size int) ([]model.Role, dto.Paging, error) {
	skip := (page - 1) * size
	qSelect := &query.QSelect{DB: r.db}
	qSelect.Table("mst_role").
		Column("id", "position").
		Limit(size).
		Offset(skip)

	rows, err := qSelect.Run()
	if err != nil {
		return nil, dto.Paging{}, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		err = rows.Scan(
			&role.Id,
			&role.Position,
		)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		roles = append(roles, role)
	}

	var totalRows int
	err = r.db.QueryRow("SELECT COUNT(*) FROM mst_role").Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	resultPagingDto := utils.Paginate(page, size, totalRows)
	return roles, resultPagingDto, nil
}

func (r *roleRepository) UpdateRole(newRole model.Role) error {
	qUpdate := &query.QUpdate{DB: r.db}
	qUpdate.Table("mst_role").
		Set("position", newRole.Position).
		Where("id", "=", newRole.Id)

	_, err := qUpdate.Run()
	if err != nil {
		return err
	}
	return nil
}

func (r *roleRepository) GetRoleById(id string) (model.Role, error) {
	qSelect := &query.QSelect{DB: r.db}
	qSelect.Table("mst_role").
		Column("id", "position").
		Where("id", "=", id)

	row := qSelect.RunRow()
	var role model.Role
	err := row.Scan(
		&role.Id,
		&role.Position,
	)

	if err != nil {
		return model.Role{}, err
	}
	return role, nil
}

func (r *roleRepository) DeleteRoleById(id string) error {
	qDelete := &query.QDelete{DB: r.db}
	qDelete.Table("mst_role").Where("id", "=", id)

	_, err := qDelete.Run()
	if err != nil {
		return err
	}
	return nil
}

func newRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) InsertRole(newRole model.Role) error {
	qInsert := query.QInsert{DB: r.db}
	qInsert.Table("mst_role").
		Column("id", "position").
		Values(newRole.Id, newRole.Position)

	_, err := qInsert.Run() // Placeholder parameter

	if err != nil {
		return err
	}
	return nil
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

// ada kebutuhan pengambilan data/query -> Query(list) /QueryRow(single Value)
// tidak kebutuhan untuk ngambil data/record dari db -> Exec() : contoh : insert, Update, Delete,dan lain lain.
