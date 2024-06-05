package usecase

import (
	//"errors"

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
	ChangeStatus(reserv_status request.ReservationStatusRequest) error
	GetListByToken(fl_reserv_get_list request.ReservationGetListFilter) ([]response.ReservationResponse, error)
}

type reservationUsecase struct {
	repo repository.ReservationRepository
}

func (self *reservationUsecase) CreateRequest(new_request request.ReservationRequest) ([]response.ReservationResponse, error) {
	new_request.Id = uuid.NewString()
	new_request.StatusId = utils.RESERV_STATUS_PENDING

	// Set Time Now fo Reservation DateTime
	currentTime := time.Now()
	new_request.ReservationDate = currentTime.Format("2006-01-02 15:04:05")

	// Set ID for reservation request detail
	for i, _ := range new_request.AdditionalFacility {
		new_request.AdditionalFacility[i].Id = uuid.NewString()
	}

	err := self.repo.CreateRequest(new_request)

	if err != nil {
		return nil, err
	}

	fl_reserv_get_list := request.ReservationGetListFilter{
		ReservationId: new_request.Id,
	}

	return self.GetListByToken(fl_reserv_get_list)
}

func (self *reservationUsecase) ChangeStatus(reserv_status request.ReservationStatusRequest) error {
	return self.repo.ChangeStatus(reserv_status)
}

func (self *reservationUsecase) GetListByToken(fl_reserv_get_list request.ReservationGetListFilter) ([]response.ReservationResponse, error) {
	return self.repo.GetListByToken(fl_reserv_get_list)
}

func NewReservationUsecase(repo repository.ReservationRepository) ReservationUsecase {
	return &reservationUsecase{
		repo: repo,
	}
}
