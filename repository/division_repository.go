package repository

import (
	"database/sql"
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/utils"
	"go-roomify/utils/query"
)

type DivisionRepository interface {
	GetAllDivisi(page, size int) ([]model.Division, dto.Paging, error)
	GetDivisiById(id string) (model.Division, error)
	CreateDivisi(division model.Division) error
	UpdateDivisi(division model.Division) error
	DeleteDivisi(id string) error
}

type divisionRepository struct {
	db *sql.DB
}

func (r *divisionRepository) GetAllDivisi(page, size int) ([]model.Division, dto.Paging, error) {
	skip := (page - 1) * size
	qselect := query.QSelect{DB: r.db}

	qselect.Table("mst_division")
	qselect.Column(
		"id",
		"name")
	qselect.Limit(size)
	qselect.Offset(skip)
	rows, err := qselect.Run()
	if err != nil {
		return nil, dto.Paging{}, err
	}
	defer rows.Close()

	var divisions []model.Division
	for rows.Next() {
		var division model.Division
		if err := rows.Scan(
			&division.Id,
			&division.Name); err != nil {
			return nil, dto.Paging{}, err
		}
		divisions = append(divisions, division)
	}
	var totalRows int
	qcount := query.QSelect{DB: r.db}
	qcount.Table("mst_division")
	qcount.Column("COUNT(id)")

	err = qcount.RunRow().Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	resultPagingDto := utils.Paginate(page, size, totalRows)
	return divisions, resultPagingDto, nil

}

func (r *divisionRepository) GetDivisiById(id string) (model.Division, error) {
	qselect := query.QSelect{DB: r.db}

	qselect.Table("mst_division")
	qselect.Column("id", "name")
	qselect.Where("id", "=", id)

	row := qselect.RunRow()
	var division model.Division
	err := row.Scan(&division.Id, &division.Name)
	if err != nil {
		return division, err
	}

	return division, nil
}

func (r *divisionRepository) CreateDivisi(division model.Division) error {
	qinsert := query.QInsert{DB: r.db}
	qinsert.Table("mst_division").Column("id", "name").Values(division.Id, division.Name)
	_, err := qinsert.Run()
	return err
}

func (r *divisionRepository) UpdateDivisi(division model.Division) error {
	qUpdate := query.QUpdate{DB: r.db}
	qUpdate.Table("mst_division").
		Set("name", division.Name).
		Where("id", "=", division.Id)
	_, err := qUpdate.Run()
	return err
}

func (r *divisionRepository) DeleteDivisi(id string) error {
	qdelete := query.QDelete{DB: r.db}
	qdelete.Table("mst_division").Where("id", "=", id)
	_, err := qdelete.Run()
	return err
}

func NewDivisionRepository(db *sql.DB) DivisionRepository {
	return &divisionRepository{db: db}
}
