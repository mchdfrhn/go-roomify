package repository

import (
	"database/sql"
	"fmt"
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/utils"
	"go-roomify/utils/query"
	"strconv"
)

type FacilityRepository interface {
	GetFacilityById(id string) (model.Facility, error)
	GetPagingFacility(page int, size int) ([]model.Facility, dto.Paging, error)
	InsertFacility(newFacility model.Facility) error
	UpdateFacility(newFacility model.Facility) error
	DeleteFacility(id string) error
}

type facilityRepository struct {
	db *sql.DB
}

func NewFacilityRepository(db *sql.DB) FacilityRepository {
	return &facilityRepository{
		db: db,
	}
}

func (f *facilityRepository) GetFacilityById(id string) (model.Facility, error) {
	var facility model.Facility
	qselect := query.QSelect{DB: f.db}
	qselect.Table("mst_facility")
	qselect.Column(
		"id",
		"name",
		"is_available",
		"is_reserveable",
		"room_id",
	)
	qselect.Where("id", "=", id)
	rows := qselect.RunRow()
	err := rows.Scan(
		&facility.Id,
		&facility.Name,
		&facility.IsAvailable,
		&facility.IsReserveable,
		&facility.RoomId,
	)

	if err != nil {
		return model.Facility{}, err
	}

	return facility, nil

}

func (f *facilityRepository) GetPagingFacility(page int, size int) ([]model.Facility, dto.Paging, error) {
	skip := (page - 1) * size
	qselect := query.QSelect{DB: f.db}
	qselect.Table("mst_facility")
	qselect.Column(
		"id",
		"name",
		"is_available",
		"is_reserveable",
		"room_id",
	)
	qselect.Limit(size)
	qselect.Offset(skip)

	rows, err := qselect.Run()
	if err != nil {
		return nil, dto.Paging{}, err
	}

	var facilities []model.Facility
	for rows.Next() {
		var facility model.Facility
		err = rows.Scan(
			&facility.Id,
			&facility.Name,
			&facility.IsAvailable,
			&facility.IsReserveable,
			&facility.RoomId,
		)
		if err != nil {
			return nil, dto.Paging{}, err
		}

		facilities = append(facilities, facility)
	}

	var totalRows int
	qcount := query.QSelect{DB: f.db}
	qcount.Table("mst_facility")
	qcount.Column("COUNT(id)")

	err = qcount.RunRow().Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	resultPagingDto := utils.Paginate(page, size, totalRows)
	return facilities, resultPagingDto, nil
}

func (f *facilityRepository) InsertFacility(newFacility model.Facility) error {
	qinsert := query.QInsert{DB: f.db}
	qinsert.Table("mst_facility")
	qinsert.Column(
		"id",
		"name",
		"is_available",
		"is_reserveable",
		"room_id",
	)
	qinsert.Values(
		newFacility.Id,
		newFacility.Name,
		strconv.FormatBool(newFacility.IsAvailable),
		strconv.FormatBool(newFacility.IsReserveable),
		newFacility.RoomId,
	)

	fmt.Println(qinsert.GetQuery())
	_, err := qinsert.Run()
	if err != nil {
		return err
	}

	return nil

}

func (f *facilityRepository) UpdateFacility(newFacility model.Facility) error {
	qupdate := query.QUpdate{DB: f.db}
	qupdate.Table("mst_facility")
	qupdate.Set("name", newFacility.Name)
	qupdate.Set("is_available", strconv.FormatBool(newFacility.IsAvailable))
	qupdate.Set("is_reserveable", strconv.FormatBool(newFacility.IsReserveable))
	qupdate.Set("room_id", newFacility.RoomId)
	qupdate.Where("id", "=", newFacility.Id)

	fmt.Println(qupdate.GetQuery())
	_, err := qupdate.Run()
	if err != nil {
		return err
	}

	return nil

}

func (f *facilityRepository) DeleteFacility(id string) error {
	qdelete := query.QDelete{DB: f.db}
	qdelete.Table("mst_facility")
	qdelete.Where("id", "=", id)
	_, err := qdelete.Run()
	if err != nil {
		return err
	}

	return nil
}
