package usecase

import (
	//"errors"
	"go-roomify/model"
	"go-roomify/model/dto/request"
	"go-roomify/repository"
	"time"

	// "fmt"

	"github.com/google/uuid"
)

type ReservationUsecase interface {
	CreateRequest(new_request model.Reservation) (model.Reservation, error)
	ChangeStatus(reserv_status request.ReservationStatus) error
}

type reservationUsecase struct {
	repo repository.ReservationRepository
}

func (self *reservationUsecase) CreateRequest(new_request model.Reservation) (model.Reservation, error) {
	new_request.Id = uuid.NewString()
	new_request.Status = 0
	//new_request.Description = ""

	// Set Time Now fo Reservation DateTime
	currentTime := time.Now()
	new_request.ReservationDate = currentTime.Format("2006-01-02 15:04:05")

	// Set ID for reservation request detail
	for i, _ := range new_request.Detail {
		new_request.Detail[i].ReservationId = new_request.Id
		new_request.Detail[i].Id = uuid.NewString()
	}

	return self.repo.CreateRequest(new_request)
}

func (self *reservationUsecase) ChangeStatus(reserv_status request.ReservationStatus) error {
	return self.repo.ChangeStatus(reserv_status)
}

func NewReservationUsecase(repo repository.ReservationRepository) ReservationUsecase {
	return &reservationUsecase{
		repo: repo,
	}
}
