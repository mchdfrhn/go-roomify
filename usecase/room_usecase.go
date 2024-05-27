package usecase

import (
	"fmt"
	"go-roomify/model"
	"go-roomify/model/dto/request"
	"go-roomify/repository"

	"github.com/google/uuid"
)

type RoomUsecase interface{
	CreateRoom( roomRequest request.RoomRequest ) ( model.Room, error )
}

type roomUsecase struct{
	repo repository.RoomRepository
}

func ( ru *roomUsecase ) CreateRoom( roomRequest request.RoomRequest ) ( model.Room, error ){

	if len( roomRequest.Name ) > 200 {
		return model.Room{}, fmt.Errorf("name max 200 char")
	}

	if len( roomRequest.RoomType ) > 200 {
		return model.Room{}, fmt.Errorf("room type max 200 char")
	}

	roomId, err := ru.repo.GetRoomIdIfExist( roomRequest.Name, roomRequest.RoomType )
	if err != nil {
		return model.Room{}, err
	}

	if roomId != "" {
		return model.Room{}, fmt.Errorf("a room with the same name and type already exists")
	}

	roomModel := model.Room{
		Id: uuid.NewString(),
		Name: roomRequest.Name,
		RoomType: roomRequest.RoomType,
		Capacity: roomRequest.Capacity,
		IsAvailable: false,
	}

	return ru.repo.CreateRoom( roomModel )

}

func NewRoomUsecase( repo repository.RoomRepository ) RoomUsecase{
	return &roomUsecase{
		repo: repo,
	}
}