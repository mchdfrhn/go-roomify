package repository

import (
	"database/sql"
	"go-roomify/model"
	"go-roomify/model/dto/request"
	"go-roomify/model/dto/response"
	"go-roomify/utils"

	//"go-roomify/model/dto"
	//"go-roomify/model/dto/request"
	//"go-roomify/utils"
	"go-roomify/utils/query"
	// "fmt"
	"errors"
)

type ReservationRepository interface {
	CreateRequest(new_request model.Reservation) (model.Reservation, error)
	ChangeStatus(reserv_status request.ReservationStatusRequest) error
	GetListByToken(fl_reserv_get_list request.ReservationGetListFilter) ([]response.ReservationResponse, error)
	GetReservationByYear(startYear string, endYear string) (*sql.Rows, error)
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
	qselect.Where("id", "=", new_request.Detail.RoomId)
	qselect.AndWhere("is_available", "=", false)

	var is_room_booked int
	err := qselect.RunRow().Scan(&is_room_booked)

	if err != nil {
		return model.Reservation{}, err
	}

	if is_room_booked > 0 {
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
		"reservation_status_id",
		"description")
	qinsert.Values(
		new_request.Id,
		new_request.UserProfileId,
		new_request.ReservationDate,
		new_request.StartDate,
		new_request.EndDate,
		new_request.Status.Id,
		new_request.Description)

	result, err := qinsert.Run()

	if err != nil {
		return model.Reservation{}, err
	}

	if a, _ := result.RowsAffected(); a == 0 {
		return model.Reservation{}, errors.New("Failed to make a request, no rows affected")
	}

	// Insert Reservation Detail
	qinsert = query.QInsert{DB: self.db}
	req_detail := new_request.Detail

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

	result, err = qinsert.Run()

	if err != nil {
		return model.Reservation{}, err
	}

	if a, _ := result.RowsAffected(); a == 0 {
		return model.Reservation{}, errors.New("Failed to make a request, no rows affected")
	}

	return new_request, nil
}

func (self *reservationRepository) ChangeStatus(reserv_status request.ReservationStatusRequest) error {
	// SET ROOM TO UNAVAILABLE
	qselect := query.QSelect{DB: self.db}
	qselect.Table("tx_reservation AS rv")
	qselect.Column("rvd.room_id")
	qselect.Join("tx_reservation_detail as rvd", "rvd.reservation_id = rv.id")
	qselect.Where("rv.id", "=", reserv_status.ReservationId)

	var room_id string

	if err := qselect.RunRow().Scan(&room_id); err != nil {
		return err
	}

	qupdate := query.QUpdate{DB: self.db}
	qupdate.Table("mst_room")
	qupdate.Set("is_available", false)
	qupdate.Where("id", "=", room_id)

	result, err := qupdate.Run()

	if err != nil {
		return err
	}

	if a, _ := result.RowsAffected(); a == 0 {
		return errors.New("Invalid Room Id")
	}

	// Set Transaction Status
	qupdate = query.QUpdate{DB: self.db}

	qupdate.Table("tx_reservation")
	qupdate.Set("reservation_status_id", reserv_status.StatusId)
	qupdate.Set("description", reserv_status.Description)
	qupdate.Where("id", "=", reserv_status.ReservationId)
	qupdate.AndWhere("reservation_status_id", "=", utils.RESERV_STATUS_PENDING)

	result, err = qupdate.Run()

	if err != nil {
		return err
	}

	if a, _ := result.RowsAffected(); a == 0 {
		return errors.New("Invalid Reservation Id")
	}

	return err
}

func (self *reservationRepository) GetListByToken(fl_reserv_get_list request.ReservationGetListFilter) ([]response.ReservationResponse, error) {
	qselect := query.QSelect{DB: self.db}

	qselect.Table("tx_reservation AS rv")
	qselect.Column(
		"rv.id",
		"rv.user_profile_id",
		"rv.reservation_date",
		"rv.start_time",
		"rv.end_time",
		"rv.description",

		"rvs.id AS status_id",
		"rvs.name AS status_name",

		"rvd.id",
		"rvd.reservation_id",
		"rvd.equipment_needed",

		"room.id",
		"room.name",
		"room.roomtype",
		"room.capacity",
		"room.is_available")
	qselect.Join("tx_reservation_status AS rvs", "rv.reservation_status_id = rvs.id")
	qselect.Join("tx_reservation_detail AS rvd", "rv.id = rvd.reservation_id")
	qselect.Join("mst_room AS room", "rvd.room_id = room.id")

	filter_status := fl_reserv_get_list.FilterStatus
	user_id := fl_reserv_get_list.UserId
	user_role := fl_reserv_get_list.UserRole

	if filter_status != "" {
		qselect.Where("rvs.id", "=", filter_status)
	}
	if user_role == "employee" {
		if filter_status != "" {
			qselect.AndWhere("rv.user_profile_id", "=", user_id)
		} else {
			qselect.Where("rv.user_profile_id", "=", user_id)
		}
	}

	rows, err := qselect.Run()

	if err != nil {
		return nil, err
	}

	var rows_reserv_response []response.ReservationResponse

	for rows.Next() {
		var r_reserv_resp response.ReservationResponse

		err := rows.Scan(
			&r_reserv_resp.Id,
			&r_reserv_resp.UserProfileId,
			&r_reserv_resp.ReservationDate,
			&r_reserv_resp.StartDate,
			&r_reserv_resp.EndDate,
			&r_reserv_resp.Description,
			&r_reserv_resp.Status.Id,
			&r_reserv_resp.Status.Name,
			&r_reserv_resp.Detail.Id,
			&r_reserv_resp.Detail.ReservationId,
			&r_reserv_resp.Detail.Equipment,
			&r_reserv_resp.Detail.Room.Id,
			&r_reserv_resp.Detail.Room.Name,
			&r_reserv_resp.Detail.Room.RoomType,
			&r_reserv_resp.Detail.Room.Capacity,
			&r_reserv_resp.Detail.Room.IsAvailable)

		if err != nil {
			return nil, err
		}

		rows_reserv_response = append(rows_reserv_response, r_reserv_resp)
	}

	rows.Close()

	return rows_reserv_response, nil
}

func (self *reservationRepository) GetReservationByYear(startYear string, endYear string) (*sql.Rows, error) {

	query := query.QSelect{DB: self.db}

	rows, err := query.Table(
		"tx_reservation AS tr",
	).Column(
		"*",
	).Join(
		"tx_reservation_detail AS trd",
		"trd.reservation_id = tr.id",
	).Join(
		"room AS r",
		"r.id = trd.room_id",
	).Join(
		"tx_reservation_status AS trs",
		"trs.id = tr.reservation_status_id",
	).Join(
		"facility AS f",
		"f.room_id = r.id",
	).Where(
		"reservation_date",
		">=",
		startYear,
	).AndWhere(
		"reservation_date",
		"<=",
		endYear,
	).Run()

	// rows, err := query.Table(
	// 	"t_tes",
	// ).Column(
	// 	"*",
	// ).Where(
	// 	"date",
	// 	">=",
	// 	startYear,
	// ).AndWhere(
	// 	"date",
	// 	"<=",
	// 	endYear,
	// ).Run()

	if err != nil {
		return nil, err
	}

	return rows, nil

}

func NewReservationRepository(db *sql.DB) ReservationRepository {
	return &reservationRepository{
		db: db,
	}
}
