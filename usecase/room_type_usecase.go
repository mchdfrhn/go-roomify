package usecase

import (
	"fmt"
	"go-roomify/model"
	"go-roomify/model/dto/request"
	"go-roomify/repository"
	"net/http"

	"github.com/google/uuid"
)

type RoomTypeUsecase interface {
	CreateRoomType(roomRequest request.RoomTypeRequest) (model.RoomType, int, error)
	GetAllRoomType() ([]model.RoomType, int, error)
	GetRoomTypeByIdOrName(idOrNameRoom string) (model.RoomType, int, error)
	UpdateRoomTypeById(updateRoomRequest model.RoomType) (model.RoomType, int, error)
	DeleteRoomTypeById(roomId string) (int, error)
}

type roomTypeUsecase struct {
	repo repository.RoomTypeRepository
}

func (ru *roomTypeUsecase) CreateRoomType(roomTypeRequest request.RoomTypeRequest) (model.RoomType, int, error) {
	if len(roomTypeRequest.Name) > 200 {
		return model.RoomType{}, http.StatusBadRequest, fmt.Errorf("name max 200 char")
	}

	roomType, err := ru.repo.GetRoomTypeByIdOrName(roomTypeRequest.Name)
	if err != nil {
		return model.RoomType{}, http.StatusInternalServerError, err
	}

	if roomType.Id != "" {
		return model.RoomType{}, http.StatusConflict, fmt.Errorf("a room with the same name already exists")
	}

	roomTypeModel := model.RoomType{
		Id:            uuid.NewString(),
		Name:          roomTypeRequest.Name,
	}

	err = ru.repo.CreateRoomType(roomTypeModel)
	if err != nil {
		return model.RoomType{}, http.StatusInternalServerError, err
	}

	return roomTypeModel, http.StatusCreated, nil
}

func (ru *roomTypeUsecase) GetAllRoomType() ([]model.RoomType, int, error) {
	allRoomType, err := ru.repo.GetAllRoomType()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return allRoomType, http.StatusOK, nil
}

func (ru *roomTypeUsecase) GetRoomTypeByIdOrName(roomTypeIdOrName string) (model.RoomType, int, error) {
	findRoomType, err := ru.repo.GetRoomTypeByIdOrName(roomTypeIdOrName)
	if err != nil {
		return model.RoomType{}, http.StatusInternalServerError, err
	}

	if findRoomType.Id == "" {
		return model.RoomType{}, http.StatusNotFound, fmt.Errorf("room type not found")
	}

	return findRoomType, http.StatusOK, nil
}

func (ru *roomTypeUsecase) UpdateRoomTypeById(updateRoom model.RoomType) (model.RoomType, int, error) {
	if len(updateRoom.Name) > 200 {
		return model.RoomType{}, http.StatusBadRequest, fmt.Errorf("name max 200 char")
	}

	_, status, err := ru.GetRoomTypeByIdOrName(updateRoom.Id)
	if err != nil {
		return model.RoomType{}, status, err
	}

	updatedRoom, err := ru.repo.UpdateRoomTypeById(updateRoom)
	if err != nil {
		return model.RoomType{}, http.StatusInternalServerError, err
	}

	return updatedRoom, http.StatusOK, nil
}

func (ru *roomTypeUsecase) DeleteRoomTypeById(roomTypeId string) (int, error) {
	_, status, err := ru.GetRoomTypeByIdOrName(roomTypeId)
	if err != nil {
		return status, err
	}

	err = ru.repo.DeleteRoomTypeById(roomTypeId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func NewRoomTypeUsecase(repo repository.RoomTypeRepository) RoomTypeUsecase {
	return &roomTypeUsecase{
		repo: repo,
	}
}
