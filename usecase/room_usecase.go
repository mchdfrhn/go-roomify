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
	GetAllRoom(paramPage string, paramSize string, paramType string) ([]any, dto.Paging, int, error)
	GetRoomByIdOrName(idOrNameRoom string) ([]response.RoomResponse, int, error)
	UpdateRoomById(updateRoomRequest request.UpdateRoomRequest) (request.UpdateRoomRequest, int, error)
	DeleteRooomById(roomId string) (int, error)
	UpdateRoomByIdIsAvailableOnly(request.RoomStatusRequest) (request.RoomStatusRequest, int, error)
	GetAvailableRoom(paramType string) ([]response.RoomResponse, int, error)
}

type roomUsecase struct {
	repo       repository.RoomRepository
	roomTypeUc RoomTypeUsecase
}

func (ru *roomUsecase) CreateRoom(roomRequest request.RoomRequest) (model.Room, int, error) {
	if len(roomRequest.Name) > 200 {
		return model.Room{}, http.StatusBadRequest, fmt.Errorf("name max 200 char")
	}

	if len(roomRequest.RoomTypeId) > 200 {
		return model.Room{}, http.StatusBadRequest, fmt.Errorf("room type max 200 char")
	}

	_, status, err := ru.roomTypeUc.GetRoomTypeByIdOrName(roomRequest.RoomTypeId)
	if err != nil {
		return model.Room{}, status, err
	}

	roomId, err := ru.repo.GetRoomIdIfExist(roomRequest.Name, roomRequest.RoomTypeId)
	if err != nil {
		return model.Room{}, http.StatusInternalServerError, fmt.Errorf("server error")
	}

	if roomId != "" {
		return model.Room{}, http.StatusConflict, fmt.Errorf("a room with the same name and type already exists")
	}

	isAvailable := true
	roomModel := model.Room{
		Id:            uuid.NewString(),
		Name:          roomRequest.Name,
		RoomTypeId:    roomRequest.RoomTypeId,
		Capacity:      roomRequest.Capacity,
		IsAvailable:   &isAvailable,
		IsReserveable: roomRequest.IsReserveable,
	}

	err = ru.repo.CreateRoom(roomModel)
	if err != nil {
		return model.Room{}, http.StatusInternalServerError, fmt.Errorf("server error")
	}

	return roomModel, http.StatusCreated, nil
}

func (ru *roomUsecase) GetAllRoom(paramPage string, paramSize string, paramType string) ([]any, dto.Paging, int, error) {
	page := 1
	size := 10

	if paramPage != "" {
		page, _ = strconv.Atoi(paramPage)
	}

	if paramSize != "" {
		size, _ = strconv.Atoi(paramSize)
	}

	if page <= 0 || size <= 0 {
		return nil, dto.Paging{}, http.StatusBadRequest, fmt.Errorf("page or size number must be a positive integer")
	}

	skip := (page - 1) * size
	allRoom, paging, err := ru.repo.GetAllRoom(page, skip, size, paramType)
	if err != nil {
		return nil, dto.Paging{}, http.StatusInternalServerError, fmt.Errorf("server error")
	}

	var castingAllRoom []any
	castingAllRoom = append(castingAllRoom, allRoom)
	return castingAllRoom, paging, http.StatusOK, nil
}

func (ru *roomUsecase) GetRoomByIdOrName(roomIdOrName string) ([]response.RoomResponse, int, error) {
	findRoom, err := ru.repo.GetRoomByIdOrName(roomIdOrName)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("server error")
	}

	if len(findRoom) == 0 {
		return nil, http.StatusNotFound, fmt.Errorf("room not found")
	}

	return findRoom, http.StatusOK, nil
}

func (ru *roomUsecase) UpdateRoomById(updateRoom request.UpdateRoomRequest) (request.UpdateRoomRequest, int, error) {
	if len(updateRoom.Name) > 200 {
		return request.UpdateRoomRequest{}, http.StatusBadRequest, fmt.Errorf("name max 200 char")
	}

	if len(updateRoom.RoomTypeId) > 200 {
		return request.UpdateRoomRequest{}, http.StatusBadRequest, fmt.Errorf("room type max 200 char")
	}

	findRoom, status, err := ru.GetRoomByIdOrName(updateRoom.Id)
	if err != nil {
		return request.UpdateRoomRequest{}, status, err
	}

	if len(findRoom) > 1 {
		return request.UpdateRoomRequest{}, http.StatusBadRequest, fmt.Errorf("by id not by name")
	}

	if findRoom[0].Id == updateRoom.Id &&
		findRoom[0].Name == updateRoom.Name &&
		findRoom[0].Capacity == updateRoom.Capacity &&
		findRoom[0].IsReserveable == *updateRoom.IsReserveable {

		return request.UpdateRoomRequest{}, http.StatusBadRequest, fmt.Errorf("no changes detected")
	}

	findIdSameRoomAndType, err := ru.repo.GetRoomIdIfExist(updateRoom.Name, updateRoom.RoomTypeId)
	if err != nil {
		return request.UpdateRoomRequest{}, http.StatusInternalServerError, fmt.Errorf("server error")
	}

	if findIdSameRoomAndType != updateRoom.Id {
		return request.UpdateRoomRequest{}, http.StatusConflict, fmt.Errorf("a room with the same name and type already exists")
	}

	updatedRoom, err := ru.repo.UpdateRoomById(updateRoom)
	if err != nil {
		return request.UpdateRoomRequest{}, http.StatusInternalServerError, fmt.Errorf("server error")
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
		return http.StatusInternalServerError, fmt.Errorf("server error")
	}

	return http.StatusOK, nil
}

func (ru *roomUsecase) UpdateRoomByIdIsAvailableOnly(updateRoom request.RoomStatusRequest) (request.RoomStatusRequest, int, error) {
	findRoom, status, err := ru.GetRoomByIdOrName(updateRoom.Id)
	if err != nil {
		return request.RoomStatusRequest{}, status, err
	}

	if len(findRoom) > 1 {
		return request.RoomStatusRequest{}, http.StatusBadRequest, fmt.Errorf("by id not by name")
	}

	err = ru.repo.UpdateRoomByIdAvailableOnly(updateRoom)
	if err != nil {
		return request.RoomStatusRequest{}, http.StatusInternalServerError, fmt.Errorf("server error")
	}

	return updateRoom, http.StatusOK, nil
}

func (ru *roomUsecase) GetAvailableRoom(paramType string) ([]response.RoomResponse, int, error) {
	availableRooms, err := ru.repo.GetAvailableRoom(paramType)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return availableRooms, http.StatusOK, nil
}

func NewRoomUsecase(repo repository.RoomRepository, roomTypeUc RoomTypeUsecase) RoomUsecase {
	return &roomUsecase{
		repo:       repo,
		roomTypeUc: roomTypeUc,
	}
}
