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
	CreateRoomType(roomTypeRequest request.RoomTypeRequest) (model.RoomType, int, error)
	GetAllRoomType() ([]model.RoomType, int, error)
	GetRoomTypeByIdOrName(idOrNameRoomType string) (model.RoomType, int, error)
	UpdateRoomTypeById(updateRoomTypeRequest model.RoomType) (model.RoomType, int, error)
	DeleteRoomTypeById(roomTypeId string) (int, error)
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
		return model.RoomType{}, http.StatusInternalServerError, fmt.Errorf("server error")
	}

	if roomType.Id != "" {
		return model.RoomType{}, http.StatusConflict, fmt.Errorf("a room with the same name already exists")
	}

	roomTypeModel := model.RoomType{
		Id:   uuid.NewString(),
		Name: roomTypeRequest.Name,
	}

	err = ru.repo.CreateRoomType(roomTypeModel)
	if err != nil {
		return model.RoomType{}, http.StatusInternalServerError, fmt.Errorf("server error")
	}

	return roomTypeModel, http.StatusCreated, nil
}

func (ru *roomTypeUsecase) GetAllRoomType() ([]model.RoomType, int, error) {
	allRoomType, err := ru.repo.GetAllRoomType()
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("server error")
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

func (ru *roomTypeUsecase) UpdateRoomTypeById(updateRoomType model.RoomType) (model.RoomType, int, error) {
	if len(updateRoomType.Name) > 200 {
		return model.RoomType{}, http.StatusBadRequest, fmt.Errorf("name max 200 char")
	}

	findRoomTypeById, status, err := ru.GetRoomTypeByIdOrName(updateRoomType.Id)
	if err != nil {
		return model.RoomType{}, status, err
	}

	if findRoomTypeById == updateRoomType {
		return model.RoomType{}, http.StatusBadRequest, fmt.Errorf("no changes detected")
	}

	findRoomTypeByName, status, err := ru.GetRoomTypeByIdOrName(updateRoomType.Name)
	if err != nil {
		if status != 404 {
			return model.RoomType{}, status, err
		}
	}

	if findRoomTypeByName.Name == updateRoomType.Name &&
		findRoomTypeByName.Id != updateRoomType.Id {
		return model.RoomType{}, http.StatusBadRequest, fmt.Errorf("name must be unique")
	}

	updatedRoomType, err := ru.repo.UpdateRoomTypeById(updateRoomType)
	if err != nil {
		return model.RoomType{}, http.StatusInternalServerError, err
	}

	return updatedRoomType, http.StatusOK, nil
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
