package repository

import (
	"database/sql"
	"go-roomify/model"
	"go-roomify/utils/query"
)

type RoomTypeRepository interface {
	CreateRoomType(roomModelType model.RoomType) error
	GetAllRoomType() ([]model.RoomType, error)
	GetRoomTypeByIdOrName(idOrNameRoomType string) (model.RoomType, error)
	UpdateRoomTypeById(updateRoomType model.RoomType) (model.RoomType, error)
	DeleteRoomTypeById(roomId string) error
}

type roomTypeRepository struct {
	db *sql.DB
}

func (rr *roomTypeRepository) CreateRoomType(roomModelType model.RoomType) error {
	query := query.QInsert{DB: rr.db}

	_, err := query.Table(
		"room_type",
	).Column(
		"id",
		"name",
	).Values(
		roomModelType.Id,
		roomModelType.Name,
	).Run()

	if err != nil {
		return err
	}

	return nil
}

func (rr *roomTypeRepository) GetAllRoomType() ([]model.RoomType, error) {
	queryAllRoom := query.QSelect{DB: rr.db}
	rows, err := queryAllRoom.Table(
		"room_type",
	).Column(
		"id",
		"name",
	).Run()
	if err != nil {
		return nil, err
	}

	var allRoomType []model.RoomType
	for rows.Next() {
		var roomType model.RoomType

		err := rows.Scan(
			&roomType.Id,
			&roomType.Name,
		)
		if err != nil {
			return nil, err
		}

		allRoomType = append(allRoomType, roomType)
	}

	return allRoomType, nil
}

func (rr *roomTypeRepository) GetRoomTypeByIdOrName(roomTypeIdOrName string) (model.RoomType, error) {
	query := query.QSelect{DB: rr.db}
	var roomType model.RoomType

	err := query.Table(
		"room_type",
	).Column(
		"id",
		"name",
	).Where(
		"id",
		"=",
		roomTypeIdOrName,
	).OrWhere(
		"name",
		"=",
		roomTypeIdOrName,
	).RunRow().Scan(
		&roomType.Id,
		&roomType.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.RoomType{}, nil
		}
		return model.RoomType{}, err
	}

	return roomType, nil
}

func (rr *roomTypeRepository) UpdateRoomTypeById(updateRoomType model.RoomType) (model.RoomType, error) {
	query := query.QUpdate{DB: rr.db}

	_, err := query.Table(
		"room_type",
	).Set(
		"name",
		updateRoomType.Name,
	).Where(
		"id",
		"=",
		updateRoomType.Id,
	).Run()

	if err != nil {
		return model.RoomType{}, err
	}

	return updateRoomType, nil
}

func (rr *roomTypeRepository) DeleteRoomTypeById(roomId string) error {
	query := query.QDelete{DB: rr.db}

	_, err := query.Table(
		"room_type",
	).Where(
		"id",
		"=",
		roomId,
	).Run()

	if err != nil {
		return err
	}

	return nil

}

func NewRoomTypeRepository(db *sql.DB) RoomTypeRepository {
	return &roomTypeRepository{
		db: db,
	}
}
