package repository

import (
	"database/sql"
	"go-roomify/model"
	"go-roomify/utils/query"
)

type RoomRepository interface{
	CreateRoom( roomModel model.Room ) error
	GetRoomIdIfExist( name string, roomtype string ) ( string, error )
	GetRoomByIdOrName( idOrNameRoom string ) ( []model.Room, error )
	UpdateRoomById( updateRoom model.Room ) ( model.Room, error )
}

type roomRepository struct{
	db *sql.DB
}

func ( rr *roomRepository ) CreateRoom( roomModel model.Room ) error {

	query := query.QInsert{ DB: rr.db }

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

func ( rr *roomRepository ) GetRoomIdIfExist( name string, roomtype string ) ( string, error ){

	query := query.QSelect{ DB: rr.db }
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
	).RunRow().Scan( &idRoom )

	if err != nil {
		return "", nil
	}
	
	return idRoom, nil
}

func ( rr *roomRepository ) GetRoomByIdOrName( roomidOrName string ) ( []model.Room, error ){

	query := query.QSelect{ DB: rr.db }
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

func ( rr *roomRepository ) UpdateRoomById( updateRoom model.Room ) ( model.Room, error ){

	query := query.QUpdate{ DB: rr.db }

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

func NewRoomRepository( db *sql.DB ) RoomRepository{
	return &roomRepository{
		db: db,
	}
}