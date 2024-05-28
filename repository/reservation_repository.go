package repository

import (
	"database/sql"
	"go-roomify/model"
	//"go-roomify/model/dto"
	//"go-roomify/model/dto/request"
	//"go-roomify/utils"
	"go-roomify/utils/query"
	// "fmt"
	"errors"
)

type ReservationRepository interface {
	CreateRequest(new_request model.Reservation) (model.Reservation, error)
	CancelRequest(id string) (error)
}

type reservationRepository struct {
	db *sql.DB
}

func (self *reservationRepository) CreateRequest(new_request model.Reservation) (model.Reservation, error) {
	//var r_reservation response.ReservationResponse

	// [verif] Cek apakah ruangan sudah dipesan atau belum
	// Sesuai dengan 
	qselect := query.QSelect{DB: self.db}
	qselect.Table("mst_room")
	qselect.Column("COUNT(*)")
	qselect.Where("", "", false)

	for _, req_detail := range new_request.Detail {
		qselect.OrWhere("id", "=", req_detail.RoomId)
		qselect.AndWhere("is_available", "=", false)
	}

	var total_booked_room int
	err := qselect.RunRow().Scan(&total_booked_room)

	if err != nil {
		return model.Reservation{}, err
	}

	if total_booked_room > 0 {
		return model.Reservation{}, errors.New("Invalid Room")
	}

	// Insert Reservation
	qinsert := query.QInsert{DB: self.db}

	qinsert.Table("tx_reservation")
	qinsert.Column(
		"id",
		"user_profile_id",
		"reservation_date",
		"start_time",
		"end_time",
		"status",
		"description")
	qinsert.Values(
		new_request.Id,
		new_request.UserProfileId,
		new_request.ReservationDate,
		new_request.StartDate,
		new_request.EndDate,
		new_request.Status,
		new_request.Description)

	result, err := qinsert.Run()

	if err != nil {
		return model.Reservation{}, err
	}

	if a, _ := result.RowsAffected(); a == 0 {
		return model.Reservation{}, errors.New("Failed to make a request, no rows affected")
	}

	for _, req_detail := range new_request.Detail {
		// Insert Reservation Detail
		qinsert := query.QInsert{DB: self.db}

		qinsert.Table("tx_reservation_detail")
		qinsert.Column(
			"id",
			"reservation_id",
			"room_id",
			"equipment_needed")
		qinsert.Values(
			req_detail.Id,
			new_request.Id,
			req_detail.RoomId,
			req_detail.Equipment)

		result, err := qinsert.Run()

		if err != nil {
			return model.Reservation{}, err
		}

		if a, _ := result.RowsAffected(); a == 0 {
			return model.Reservation{}, errors.New("Failed to make a request, no rows affected")
		}
	}

	return new_request, nil
}

func (self *reservationRepository) CancelRequest(id string) (error) {
	qupdate := query.QUpdate{DB: self.db}

	qupdate.Table("tx_reservation")
	qupdate.Set("status", 2)
	qupdate.Where("id", "=", id)
	qupdate.AndWhere("status", "=", 0)

	result, err := qupdate.Run()

	if err != nil {
		return err
	}

	if a, _ := result.RowsAffected(); a == 0 {
		return errors.New("Invalid Id")
	}

	return err
}

func NewReservationRepository(db *sql.DB) (ReservationRepository) {
	return &reservationRepository {
		db: db,
	}
}