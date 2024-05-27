package repository

import (
	"database/sql"
	"go-roomify/model"
	"go-roomify/utils/query"
)

type RoomRepository interface{
	CreateRoom( roomModel model.Room ) ( model.Room, error )
	GetRoomIdIfExist( name string, roomtype string ) ( string, error )
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

func NewRoomRepository( db *sql.DB ) RoomRepository{
	return &roomRepository{
		db: db,
	}
}