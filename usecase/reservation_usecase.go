package usecase

import (
	//"errors"

	"errors"
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/model/dto/request"
	"go-roomify/model/dto/response"
	"go-roomify/repository"
	"go-roomify/utils"
	"time"

	// "fmt"

	"github.com/google/uuid"
)

type ReservationUsecase interface {
	CreateRequest(new_request request.ReservationRequest) ([]response.ReservationResponse, error)
	ChangeStatus(reserv_status request.ReservationStatusRequest) ([]response.ReservationResponse, error)
	GetListByToken(fl_reserv_get_list request.ReservationGetListFilter) ([]response.ReservationResponse, dto.Paging, error)
}

type reservationUsecase struct {
	repo          repository.ReservationRepository
	repo_room     repository.RoomRepository
	repo_facility repository.FacilityRepository
}

func (self *reservationUsecase) CreateRequest(new_request request.ReservationRequest) ([]response.ReservationResponse, error) {
	// Verify Room Availability
	rows_room, err := self.repo_room.GetRoomByIdOrName(new_request.RoomId)

	if err != nil {
		return nil, err
	}

	if len(rows_room) != 1 {
		return nil, errors.New("Invalid Room Id")
	}

	row_room := rows_room[0]

	if !row_room.IsAvailable || !row_room.IsReserveable {
		return nil, errors.New("Room Not Available or Not Reserveable")
	}

	// Verify Faciliry Availability
	var rows_facility model.Facility

	for _, detail := range new_request.AdditionalFacility {
		rows_facility, err = self.repo_facility.GetFacilityById(detail.FacilityId)

		if err != nil {
			return nil, err
		}

		if !rows_facility.IsAvailable || !rows_facility.IsReserveable {
			return nil, errors.New("Facility Not Available or Not Reserveable")
		}
	}

	new_request.Id = uuid.NewString()
	new_request.StatusId = utils.RESERV_STATUS_PENDING

	// Set Time Now fo Reservation DateTime
	currentTime := time.Now()
	new_request.ReservationDate = currentTime.Format("2006-01-02 15:04:05")

	// Set ID for reservation request detail
	for i, _ := range new_request.AdditionalFacility {
		new_request.AdditionalFacility[i].Id = uuid.NewString()
	}

	err = self.repo.CreateRequest(new_request)

	if err != nil {
		return nil, err
	}

	fl_reserv_get_list := request.ReservationGetListFilter{
		ReservationId: new_request.Id,
	}

	row_reserv, _, err := self.GetListByToken(fl_reserv_get_list)

	return row_reserv, err
}

func (self *reservationUsecase) ChangeStatus(reserv_status request.ReservationStatusRequest) ([]response.ReservationResponse, error) {
	// Verify Room Availability
	err := self.repo.ChangeStatus(reserv_status)

	if err != nil {
		return nil, err
	}

	fl_reserv_get_list := request.ReservationGetListFilter{
		ReservationId: reserv_status.ReservationId,
	}

	row_reserv, _, err := self.GetListByToken(fl_reserv_get_list)

	return row_reserv, err
}

func (self *reservationUsecase) GetListByToken(fl_reserv_get_list request.ReservationGetListFilter) ([]response.ReservationResponse, dto.Paging, error) {
	return self.repo.GetListByToken(fl_reserv_get_list)
}

func NewReservationUsecase(
	repo repository.ReservationRepository,
	repo_room repository.RoomRepository,
	repo_facility repository.FacilityRepository) ReservationUsecase {
	return &reservationUsecase{
		repo:          repo,
		repo_room:     repo_room,
		repo_facility: repo_facility,
	}
}
