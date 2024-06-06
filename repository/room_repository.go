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
	GetRoomIdIfExist(name string, roomtypeId string) (string, error)
	GetAllRoom(page int, skip int, size int, paramType string) ([]response.RoomResponse, dto.Paging, error)
	GetRoomByIdOrName(idOrNameRoom string) ([]response.RoomResponse, error)
	UpdateRoomById(updateRoom model.Room) (model.Room, error)
	DeleteRoomById(roomId string) error
	GetAvailableRoom(paramType string) ([]response.RoomResponse, error)
}

type roomRepository struct {
	db *sql.DB
}

func (rr *roomRepository) CreateRoom(roomModel model.Room) error {
	query := query.QInsert{DB: rr.db}

	_, err := query.Table(
		"mst_room",
	).Column(
		"id", "name", "room_type_id", "capacity", "is_available", "is_reserveable",
	).Values(
		roomModel.Id, roomModel.Name, roomModel.RoomTypeId, roomModel.Capacity,
		roomModel.IsAvailable, roomModel.IsReserveable,
	).Run()

	if err != nil {
		return err
	}

	return nil
}

func (rr *roomRepository) GetRoomIdIfExist(name string, roomtypeId string) (string, error) {
	query := query.QSelect{DB: rr.db}
	var idRoom string

	err := query.Table(
		"mst_room",
	).Column(
		"id",
	).Where(
		"name", "=", name,
	).AndWhere(
		"room_type_id", "=", roomtypeId,
	).RunRow().Scan(
		&idRoom,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return idRoom, nil
}

func (rr *roomRepository) GetAllRoom(page int, skip int, size int, paramType string) ([]response.RoomResponse, dto.Paging, error) {
	queryAllRoom := query.QSelect{DB: rr.db}

	subQueryTable := fmt.Sprintf(`
		(
			SELECT r.*, res.start_time, res.end_time
			FROM mst_room AS r
			LEFT JOIN (
				SELECT room_id, start_time, end_time
				FROM tx_reservation
				ORDER BY start_time DESC
				LIMIT 1
			) AS res ON r.id = res.room_id
			LIMIT %d OFFSET %d
		) AS r
	`, size, skip)

	getRoom := queryAllRoom.Table(
		subQueryTable,
	).Column(
		"r.id",
		"r.name",
		"rt.name AS room_type",
		"r.capacity",
		"r.is_available",
		"COALESCE(r.start_time, '1945-08-17')",
		"COALESCE(r.end_time, '1945-08-17')",
		"r.is_reserveable",
		"COALESCE(f.id, 'null')",
		"COALESCE(f.name, 'null')",
		"COALESCE(f.room_id, 'null')",
	).LeftJoin(
		"mst_facility AS f", "f.room_id = r.id",
	).Join(
		"room_type AS rt", "rt.id = r.room_type_id",
	)

	if paramType != "" {
		getRoom.Where(
			"rt.name", "=", paramType,
		)
	}

	rows, err := getRoom.OrderBy(
		"r.capacity", "ASC",
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
		"rt.name AS room_type",
		"r.capacity",
		"r.is_available",
		"COALESCE(res.start_time, '1945-08-17')",
		"COALESCE(res.end_time, '1945-08-17')",
		"r.is_reserveable",
		"COALESCE(f.id, 'null')",
		"COALESCE(f.name, 'null')",
		"COALESCE(f.room_id, 'null')",
	).LeftJoin(
		"mst_facility AS f", "f.room_id = r.id",
	).Join(
		"room_type As rt", "rt.id = r.room_type_id",
	).LeftJoin(
		`
		( 
			SELECT room_id, start_time, end_time
			FROM tx_reservation
			ORDER BY start_time DESC
			LIMIT 1 
		) AS res
		`,
		"r.id = res.room_id",
	).Where(
		"r.id", "=", roomidOrName,
	).OrWhere(
		"r.name", "=", roomidOrName,
	).OrderBy(
		"r.capacity", "ASC",
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
		"room_type_id",
		updateRoom.RoomTypeId,
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
		"id","=", updateRoom.Id,
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
		"id", "=", roomId,
	).Run()

	if err != nil {
		return err
	}

	return nil

}

func (rr *roomRepository) GetAvailableRoom(paramType string) ([]response.RoomResponse, error) {
	query := query.QSelect{DB: rr.db}

	getRoom := query.Table(
		"mst_room AS r",
	).Column(
		"r.id",
		"r.name",
		"rt.name AS room_type",
		"r.capacity",
		"r.is_available",
		"COALESCE(res.start_time, '1945-08-17')",
		"COALESCE(res.end_time, '1945-08-17')",
		"r.is_reserveable",
		"COALESCE(f.id, 'null')",
		"COALESCE(f.name, 'null')",
		"COALESCE(f.room_id, 'null')",
	).LeftJoin(
		"mst_facility AS f", "f.room_id = r.id",
	).Join(
		"room_type AS rt", "rt.id = r.room_type_id",
	).LeftJoin(
		`
		( 
			SELECT room_id, start_time, end_time
			FROM tx_reservation
			ORDER BY start_time DESC
			LIMIT 1 
		) AS res
		`,
		"r.id = res.room_id",
	).Where(
		"r.is_available", "=", true,
	)

	if paramType != "" {
		getRoom.AndWhere(
			"rt.name", "=", paramType,
		)
	}

	rows, err := getRoom.OrderBy(
		"r.capacity", "ASC",
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
		var startTime string
		var endTime string

		err := rows.Scan(
			&currentRoom.Id,
			&currentRoom.Name,
			&currentRoom.RoomTypeId,
			&currentRoom.Capacity,
			&currentRoom.IsAvailable,
			&startTime,
			&endTime,
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

			if startTime == "1945-08-17" || *currentRoom.IsAvailable {
				startTime = ""
				endTime = ""
			}

			roomResponse = response.RoomResponse{
				Id:            currentRoom.Id,
				Name:          currentRoom.Name,
				RoomType:      currentRoom.RoomTypeId,
				Capacity:      currentRoom.Capacity,
				IsAvailable:   *currentRoom.IsAvailable,
				IsReserveable: *currentRoom.IsReserveable,
				StartTime:     startTime,
				EndTIme:       endTime,
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
