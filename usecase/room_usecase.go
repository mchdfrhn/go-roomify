package usecase

import (
	"fmt"
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/model/dto/request"
	"go-roomify/model/dto/response"
	"go-roomify/repository"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type RoomUsecase interface {
	CreateRoom(roomRequest request.RoomRequest) (model.Room, int, error)
	GetAllRoom(paramPage string, paramSize string) ([]any, dto.Paging, int, error)
	GetRoomByIdOrName(idOrNameRoom string) ([]response.RoomResponse, int, error)
	UpdateRoomById(updateRoomRequest model.Room) (model.Room, int, error)
	DeleteRooomById(roomId string) (int, error)
	UpdateRoomByIdIsAvailableOnly(roomId string, isAvailable bool) (model.Room, int, error)
	GetAvailableRoom() ([]response.RoomResponse, int, error)
}

type roomUsecase struct {
	repo repository.RoomRepository
}

func (ru *roomUsecase) CreateRoom(roomRequest request.RoomRequest) (model.Room, int, error) {
	if len(roomRequest.Name) > 200 {
		return model.Room{}, http.StatusBadRequest, fmt.Errorf("name max 200 char")
	}

	if len(roomRequest.RoomType) > 200 {
		return model.Room{}, http.StatusBadRequest, fmt.Errorf("room type max 200 char")
	}

	roomId, err := ru.repo.GetRoomIdIfExist(roomRequest.Name, roomRequest.RoomType)
	if err != nil {
		return model.Room{}, http.StatusInternalServerError, err
	}

	if roomId != "" {
		return model.Room{}, http.StatusConflict, fmt.Errorf("a room with the same name and type already exists")
	}

	roomModel := model.Room{
		Id:          uuid.NewString(),
		Name:        roomRequest.Name,
		RoomType:    roomRequest.RoomType,
		Capacity:    roomRequest.Capacity,
		IsAvailable: true,
		IsReserveable: true,
	}

	err = ru.repo.CreateRoom(roomModel)
	if err != nil {
		return model.Room{}, http.StatusInternalServerError, err
	}

	return roomModel, http.StatusCreated, nil
}

func (ru *roomUsecase) GetAllRoom(paramPage string, paramSize string) ([]any, dto.Paging, int, error) {
	page := 1
	size := 10

	if paramPage != "" {
		castingPage, err := strconv.Atoi(paramPage)
		if err != nil {
			return nil, dto.Paging{}, http.StatusBadRequest, err
		}

		page = castingPage
	}

	if paramSize != "" {
		castingSize, err := strconv.Atoi(paramSize)
		if err != nil {
			return nil, dto.Paging{}, http.StatusBadRequest, err
		}

		size = castingSize
	}

	if page <= 0 || size <= 0 {
		return nil, dto.Paging{}, http.StatusBadRequest, fmt.Errorf("page or size number must be a positive integer")
	}

	skip := (page - 1) * size
	allRoom, paging, err := ru.repo.GetAllRoom(page, skip, size)
	if err != nil {
		return nil, dto.Paging{}, http.StatusInternalServerError, err
	}

	var castingAllRoom []any
	castingAllRoom = append(castingAllRoom, allRoom)
	return castingAllRoom, paging, http.StatusOK, nil
}

func (ru *roomUsecase) GetRoomByIdOrName(roomIdOrName string) ([]response.RoomResponse, int, error) {
	findRoom, err := ru.repo.GetRoomByIdOrName(roomIdOrName)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(findRoom) == 0 {
		return nil, http.StatusNotFound, fmt.Errorf("room not found")
	}

	return findRoom, http.StatusOK, nil
}

func (ru *roomUsecase) UpdateRoomById(updateRoom model.Room) (model.Room, int, error) {
	if len(updateRoom.Name) > 200 {
		return model.Room{}, http.StatusBadRequest, fmt.Errorf("name max 200 char")
	}

	if len(updateRoom.RoomType) > 200 {
		return model.Room{}, http.StatusBadRequest, fmt.Errorf("room type max 200 char")
	}

	_, status, err := ru.GetRoomByIdOrName(updateRoom.Id)
	if err != nil {
		return model.Room{}, status, err
	}

	updatedRoom, err := ru.repo.UpdateRoomById(updateRoom)
	if err != nil {
		return model.Room{}, http.StatusInternalServerError, err
	}

	return updatedRoom, http.StatusOK, nil
}

func (ru *roomUsecase) DeleteRooomById(roomId string) (int, error) {
	_, status, err := ru.GetRoomByIdOrName(roomId)
	if err != nil {
		return status, err
	}

	err = ru.repo.DeleteRoomById(roomId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (ru *roomUsecase) UpdateRoomByIdIsAvailableOnly(roomId string, isAvailable bool) (model.Room, int, error) {
	findRoom, status, err := ru.GetRoomByIdOrName(roomId)
	if err != nil {
		return model.Room{}, status, err
	}

	if len(findRoom) > 1 {
		return model.Room{}, http.StatusBadRequest, fmt.Errorf("by id not by name")
	}

	updateRoom := model.Room{
		Id: findRoom[0].Id,
		Name: findRoom[0].Name,
		Capacity: findRoom[0].Capacity,
		IsAvailable: isAvailable,
		IsReserveable: findRoom[0].IsReserveable,
	}
	updatedRoom, status, err := ru.UpdateRoomById(updateRoom)
	if err != nil {
		return model.Room{}, status, err
	}

	return updatedRoom, http.StatusOK, nil
}

func (ru *roomUsecase) GetAvailableRoom() ([]response.RoomResponse, int, error) {
	availableRooms, err := ru.repo.GetAvailableRoom()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return availableRooms, http.StatusOK, nil
}

func NewRoomUsecase(repo repository.RoomRepository) RoomUsecase {
	return &roomUsecase{
		repo: repo,
	}
}
