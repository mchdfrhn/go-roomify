package repository

import (
	"database/sql"
	"fmt"
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/model/dto/response"
	"go-roomify/utils"
	"go-roomify/utils/query"
)

type RoomRepository interface {
	CreateRoom(roomModel model.Room) error
	GetRoomIdIfExist(name string, roomtype string) (string, error)
	GetAllRoom(page int, skip int, size int) ([]response.RoomResponse, dto.Paging, error)
	GetRoomByIdOrName(idOrNameRoom string) ([]response.RoomResponse, error)
	UpdateRoomById(updateRoom model.Room) (model.Room, error)
	DeleteRoomById(roomId string) error
	GetAvailableRoom() ([]response.RoomResponse, error)
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
		"is_reserveable",
	).Values(
		roomModel.Id,
		roomModel.Name,
		roomModel.RoomType,
		roomModel.Capacity,
		roomModel.IsAvailable,
		roomModel.IsReserveable,
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

func (rr *roomRepository) GetAllRoom(page int, skip int, size int) ([]response.RoomResponse, dto.Paging, error) {
	queryAllRoom := query.QSelect{DB: rr.db}

	subQueryTable := fmt.Sprintf(`
		(
			SELECT * 
			FROM mst_room 
			LIMIT %d OFFSET %d
		) AS r 
	`, size, skip)

	rows, err := queryAllRoom.Table(
		subQueryTable,
	).Column(
		"r.id",
		"r.name",
		"r.roomtype",
		"r.capacity",
		"r.is_available",
		"r.is_reserveable",
		"COALESCE(f.id, 'null')",
		"COALESCE(f.name, 'null')",
		"COALESCE(f.room_id, 'null')",
	).LeftJoin(
		"mst_facility AS f",
		"f.room_id = r.id",
	).Run()

	if err != nil {
		return nil, dto.Paging{}, err
	}

	responseData, err := rr.scanRoomAndFacility(rows)
	if err != nil {
		return nil, dto.Paging{}, err
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
	return responseData, resultPagingDto, nil
}

func (rr *roomRepository) GetRoomByIdOrName(roomidOrName string) ([]response.RoomResponse, error) {
	query := query.QSelect{DB: rr.db}

	rows, err := query.Table(
		"mst_room AS r",
	).Column(
		"r.id",
		"r.name",
		"r.roomtype",
		"r.capacity",
		"r.is_available",
		"r.is_reserveable",
		"COALESCE(f.id, 'null')",
		"COALESCE(f.name, 'null')",
		"COALESCE(f.room_id, 'null')",
	).LeftJoin(
		"mst_facility AS f",
		"f.room_id = r.id",
	).Where(
		"r.id",
		"=",
		roomidOrName,
	).OrWhere(
		"r.name",
		"=",
		roomidOrName,
	).Run()

	if err != nil {
		return nil, err
	}

	responseData, err := rr.scanRoomAndFacility(rows)
	if err != nil {
		return nil, err
	}

	return responseData, nil
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
	).Set(
		"is_reserveable",
		updateRoom.IsReserveable,
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

func (rr *roomRepository) GetAvailableRoom() ([]response.RoomResponse, error) {
	query := query.QSelect{DB: rr.db}

	rows, err := query.Table(
		"mst_room AS r",
	).Column(
		"r.id",
		"r.name",
		"r.roomtype",
		"r.capacity",
		"r.is_available",
		"r.is_reserveable",
		"COALESCE(f.id, 'null')",
		"COALESCE(f.name, 'null')",
		"COALESCE(f.room_id, 'null')",
	).LeftJoin(
		"mst_facility AS f",
		"f.room_id = r.id",
	).Where(
		"r.is_available",
		"=",
		"true",
	).Run()

	if err != nil {
		return nil, err
	}

	responseData, err := rr.scanRoomAndFacility(rows)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func (rr *roomRepository) scanRoomAndFacility(rows *sql.Rows) ([]response.RoomResponse, error) {
	var responseData []response.RoomResponse
	var roomResponse response.RoomResponse

	for rows.Next() {
		var currentRoom model.Room
		var currentFacility response.FacilityForRoomResponse

		err := rows.Scan(
			&currentRoom.Id,
			&currentRoom.Name,
			&currentRoom.RoomType,
			&currentRoom.Capacity,
			&currentRoom.IsAvailable,
			&currentRoom.IsReserveable,
			&currentFacility.Id,
			&currentFacility.Name,
			&currentFacility.RoomId,
		)
		if err != nil {
			return nil, err
		}

		if currentRoom.Id != roomResponse.Id {
			if roomResponse.Id != "" {
				responseData = append(responseData, roomResponse)
			}

			roomResponse = response.RoomResponse{
				Id:            currentRoom.Id,
				Name:          currentRoom.Name,
				RoomType:      currentRoom.RoomType,
				Capacity:      currentRoom.Capacity,
				IsAvailable:   *currentRoom.IsAvailable,
				IsReserveable: *currentRoom.IsReserveable,
			}
		}

		if currentFacility.Id != "null" {
			roomResponse.Facilities = append(roomResponse.Facilities, currentFacility)
		}
	}

	if roomResponse.Id != "" {
		responseData = append(responseData, roomResponse)
	}

	return responseData, nil
}

func NewRoomRepository(db *sql.DB) RoomRepository {
	return &roomRepository{
		db: db,
	}
}
