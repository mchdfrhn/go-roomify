package repository

import (
	"database/sql"
	"go-roomify/model"
	"go-roomify/utils/query"
)

type RoomRepository interface{
	CreateRoom( roomModel model.Room ) ( model.Room, error )
	GetRoomIdIfExist( name string, roomtype string ) ( string, error )
	GetRoomByIdOrName( idOrNameRoom string ) ( model.Room, error )
}

type roomRepository struct{
	db *sql.DB
}

func ( rr *roomRepository ) CreateRoom( roomModel model.Room ) ( model.Room, error ){

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
		return model.Room{}, err
	}
	
	return roomModel, nil
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

func ( rr *roomRepository ) GetRoomByIdOrName( roomidOrName string ) ( model.Room, error ){

	query := query.QSelect{ DB: rr.db }
	var findRoom model.Room
	
	err := query.Table( 
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
	).RunRow().Scan( 
		&findRoom.Id,
		&findRoom.Name,
		&findRoom.RoomType,
		&findRoom.Capacity,
		&findRoom.IsAvailable,
	)

	if err != nil {
		return model.Room{}, nil
	}
	
	return findRoom, nil
}

func NewRoomRepository( db *sql.DB ) RoomRepository{
	return &roomRepository{
		db: db,
	}
}