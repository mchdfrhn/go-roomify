package usecase

import (
	"fmt"
	"go-roomify/model"
	"go-roomify/model/dto/request"
	"go-roomify/repository"
	"net/http"

	"github.com/google/uuid"
)

type RoomUsecase interface{
	CreateRoom( roomRequest request.RoomRequest ) ( model.Room, int, error )
	GetRoomByIdOrName( idOrNameRoom string ) ( []model.Room, int, error )
}

type roomUsecase struct{
	repo repository.RoomRepository
}

func ( ru *roomUsecase ) CreateRoom( roomRequest request.RoomRequest ) ( model.Room, int, error ){

	if len( roomRequest.Name ) > 200 {
		return model.Room{}, http.StatusBadRequest, fmt.Errorf("name max 200 char")
	}

	if len( roomRequest.RoomType ) > 200 {
		return model.Room{}, http.StatusBadRequest, fmt.Errorf("room type max 200 char")
	}

	roomId, err := ru.repo.GetRoomIdIfExist( roomRequest.Name, roomRequest.RoomType )
	if err != nil {
		return model.Room{}, http.StatusInternalServerError, err
	}

	if roomId != "" {
		return model.Room{}, http.StatusConflict, fmt.Errorf("a room with the same name and type already exists")
	}

	roomModel := model.Room{
		Id: uuid.NewString(),
		Name: roomRequest.Name,
		RoomType: roomRequest.RoomType,
		Capacity: roomRequest.Capacity,
		IsAvailable: true,
	}

	err = ru.repo.CreateRoom( roomModel )
	if err != nil {
		return model.Room{}, http.StatusInternalServerError, err
	}

	return roomModel, http.StatusCreated, nil
}

func ( ru *roomUsecase ) GetRoomByIdOrName( roomIdOrName string ) ( []model.Room, int, error ){

	findRoom, err := ru.repo.GetRoomByIdOrName( roomIdOrName )
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if findRoom[0].Id == "" {
		return nil, http.StatusNotFound, fmt.Errorf("room not found")
	}

	return findRoom, http.StatusOK, nil

}


func NewRoomUsecase( repo repository.RoomRepository ) RoomUsecase{
	return &roomUsecase{
		repo: repo,
	}
}