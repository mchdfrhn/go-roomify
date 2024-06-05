package repository

import (
	"database/sql"
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/utils"
	"go-roomify/utils/query"
)

type RoomRepository interface {
	CreateRoom(roomModel model.Room) error
	GetRoomIdIfExist(name string, roomtype string) (string, error)
	GetAllRoom(page int, size int) ([]any, dto.Paging, error)
	GetRoomByIdOrName(idOrNameRoom string) ([]model.Room, error)
	UpdateRoomById(updateRoom model.Room) (model.Room, error)
	DeleteRoomById(roomId string) error
}

type roomRepository struct {
	db *sql.DB
}

func (rr *roomRepository) CreateRoom(roomModel model.Room) error {

	query := query.QInsert{DB: rr.db}

	_, err := query.Table(
		"mst_room",
	).Column(
		"id",
		"name",
		"roomtype",
		"capacity",
		"is_available",
	).Values(
		roomModel.Id,
		roomModel.Name,
		roomModel.RoomType,
		roomModel.Capacity,
		roomModel.IsAvailable,
	).Run()

	if err != nil {
		return err
	}

	return nil
}

func (rr *roomRepository) GetRoomIdIfExist(name string, roomtype string) (string, error) {

	query := query.QSelect{DB: rr.db}
	var idRoom string

	err := query.Table(
		"mst_room",
	).Column(
		"id",
	).Where(
		"name",
		"=",
		name,
	).AndWhere(
		"roomtype",
		"=",
		roomtype,
	).RunRow().Scan(&idRoom)

	if err != nil {
		return "", nil
	}

	return idRoom, nil
}

func (rr *roomRepository) GetAllRoom(page int, size int) ([]any, dto.Paging, error) {

	skip := (page - 1) * size
	var allRoom []any
	queryAllRoom := query.QSelect{DB: rr.db}

	rows, err := queryAllRoom.Table(
		"mst_room",
	).Column(
		"id",
		"name",
		"roomtype",
		"capacity",
		"is_available",
		"is_reserveable",
	).Limit(
		size,
	).Offset(
		skip,
	).Run()

	if err != nil {
		return nil, dto.Paging{}, err
	}

	for rows.Next() {
		var dummyRoom model.Room
		err = rows.Scan(
			&dummyRoom.Id,
			&dummyRoom.Name,
			&dummyRoom.RoomType,
			&dummyRoom.Capacity,
			&dummyRoom.IsAvailable,
			&dummyRoom.IsReserveable,
		)

		if err != nil {
			return nil, dto.Paging{}, err
		}

		allRoom = append(allRoom, dummyRoom)
	}

	var totalRows int
	queryCount := query.QSelect{DB: rr.db}
	err = queryCount.Table(
		"mst_room",
	).Column(
		"COUNT(id)",
	).RunRow().Scan(&totalRows)

	if err != nil {
		return nil, dto.Paging{}, err
	}

	resultPagingDto := utils.Paginate(page, size, totalRows)
	return allRoom, resultPagingDto, nil

}

func (rr *roomRepository) GetRoomByIdOrName(roomidOrName string) ([]model.Room, error) {

	query := query.QSelect{DB: rr.db}
	var findRoom []model.Room

	rows, err := query.Table(
		"mst_room",
	).Column(
		"id",
		"name",
		"roomtype",
		"capacity",
		"is_available",
	).Where(
		"id",
		"=",
		roomidOrName,
	).OrWhere(
		"name",
		"=",
		roomidOrName,
	).Run()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var dummyRoom model.Room
		err = rows.Scan(
			&dummyRoom.Id,
			&dummyRoom.Name,
			&dummyRoom.RoomType,
			&dummyRoom.Capacity,
			&dummyRoom.IsAvailable,
		)

		if err != nil {
			return nil, err
		}

		findRoom = append(findRoom, dummyRoom)
	}

	return findRoom, nil
}

func (rr *roomRepository) UpdateRoomById(updateRoom model.Room) (model.Room, error) {

	query := query.QUpdate{DB: rr.db}

	_, err := query.Table(
		"mst_room",
	).Set(
		"name",
		updateRoom.Name,
	).Set(
		"roomtype",
		updateRoom.RoomType,
	).Set(
		"capacity",
		updateRoom.Capacity,
	).Set(
		"is_available",
		updateRoom.IsAvailable,
	).Where(
		"id",
		"=",
		updateRoom.Id,
	).Run()

	if err != nil {
		return model.Room{}, err
	}

	return updateRoom, nil
}

func (rr *roomRepository) DeleteRoomById(roomId string) error {

	query := query.QDelete{DB: rr.db}

	_, err := query.Table(
		"mst_room",
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

func NewRoomRepository(db *sql.DB) RoomRepository {
	return &roomRepository{
		db: db,
	}
}
