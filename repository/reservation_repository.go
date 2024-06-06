package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-roomify/model/dto/request"
	"go-roomify/model/dto/response"
	"go-roomify/utils"
	"go-roomify/utils/validation"

	//"go-roomify/model/dto"
	//"go-roomify/model/dto/request"
	//"go-roomify/utils"
	"go-roomify/utils/query"
	// "fmt"
	"errors"
)

type ReservationRepository interface {
	CreateRequest(new_request request.ReservationRequest) error
	ChangeStatus(reserv_status request.ReservationStatusRequest) error
	GetListByToken(fl_reserv_get_list request.ReservationGetListFilter) ([]response.ReservationResponse, error)
}

type reservationRepository struct {
	db *sql.DB
}

func (self *reservationRepository) CreateRequest(new_request request.ReservationRequest) error {
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
		"room_id",
		"request_message")
	qinsert.Values(
		new_request.Id,
		new_request.UserProfileId,
		new_request.ReservationDate,
		new_request.StartDate,
		new_request.EndDate,
		new_request.StatusId,
		new_request.RoomId,
		new_request.RequestMessage)

	result, err := qinsert.Run()

	if err != nil {
		return err
	}

	if a, _ := result.RowsAffected(); a == 0 {
		return errors.New("Failed to make new Room Reservation Request")
	}

	// Insert Reservation Detail
	qinsert = query.QInsert{DB: self.db}
	//req_detail := new_request.Detail

	qinsert.Table("tx_reservation_detail")
	qinsert.Column(
		"id",
		"reservation_id",
		"facility_id")

	for _, additional_facility := range new_request.AdditionalFacility {
		qinsert.Values(
			additional_facility.Id,
			new_request.Id,
			additional_facility.FacilityId)
	}

	result, err = qinsert.Run()

	if err != nil {
		return err
	}

	if a, _ := result.RowsAffected(); a == 0 {
		return errors.New("Failed to make new Room Reservation Request")
	}

	return nil
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
	filter_status := fl_reserv_get_list.FilterStatus
	filter_start_date := fl_reserv_get_list.FilterStartDate
	filter_end_date := fl_reserv_get_list.FilterEndDate
	filter_room_id := fl_reserv_get_list.FilterRoomId
	resrv_id := fl_reserv_get_list.ReservationId
	user_id := fl_reserv_get_list.UserId
	user_role := fl_reserv_get_list.UserRole

	// Query SQL
	var querySQL = `
	SELECT JSON_AGG(reservation_info) AS all_reservations
	FROM (
		SELECT JSON_BUILD_OBJECT(
			'id', tx_reservation.id,
			'user_profile', (
				SELECT JSON_BUILD_OBJECT(
					'id', mst_user_profile.id,
					'full_name', mst_user_profile.full_name,
					'division', mst_division.name,
					'role', mst_role.position,
					'address', mst_user_profile.address,
					'phone_number', mst_user_profile.phone_number
				)
				FROM mst_user_profile
				LEFT JOIN mst_division ON mst_user_profile.division_id = mst_division.id
				LEFT JOIN mst_role ON mst_user_profile.role_id = mst_role.id
				WHERE mst_user_profile.id = tx_reservation.user_profile_id
			),
			'reservation_date', tx_reservation.reservation_date,
			'start_time', tx_reservation.start_time,
			'end_time', tx_reservation.end_time,
			'reservation_status', (
				SELECT JSON_BUILD_OBJECT(
					'id', tx_reservation_status.id,
					'name', tx_reservation_status.name
				)
				FROM tx_reservation_status
				WHERE tx_reservation_status.id = tx_reservation.reservation_status_id
			),
			'room', (
				SELECT JSON_BUILD_OBJECT(
					'id', mst_room.id,
					'name', mst_room.name,
					'roomtype', mst_room.roomtype,
					'capacity', mst_room.capacity,
					'is_available', mst_room.is_available,
					'is_reserveable', mst_room.is_reserveable,
					'facility', (
						SELECT COALESCE(
							JSON_AGG(
								JSON_BUILD_OBJECT(
									'id', mst_facility.id,
									'name', mst_facility.name,
									'is_available', mst_facility.is_available,
									'is_reserveable', mst_facility.is_reserveable,
									'room_id', mst_facility.room_id
								)
							), '[]'::JSON
						)
						FROM mst_facility
						WHERE mst_facility.room_id = mst_room.id
					)
				)
				FROM mst_room
				WHERE mst_room.id = tx_reservation.room_id
			),
			'request_message', tx_reservation.request_message,
			'response_message', tx_reservation.response_message,
			'additional_facility', (
				SELECT COALESCE(
					JSON_AGG(
						JSON_BUILD_OBJECT(
							'id', mst_facility.id,
							'name', mst_facility.name,
							'is_available', mst_facility.is_available,
							'is_reserveable', mst_facility.is_reserveable
						)
					), '[]'::JSON
				)
				FROM tx_reservation_detail
				LEFT JOIN mst_facility ON tx_reservation_detail.facility_id = mst_facility.id
				WHERE tx_reservation_detail.reservation_id = tx_reservation.id
			)
		) AS reservation_info
		FROM tx_reservation
		LEFT JOIN tx_reservation_status ON tx_reservation_status.id = tx_reservation.reservation_status_id
		LEFT JOIN mst_room ON mst_room.id = tx_reservation.room_id
		WHERE true `
	if user_role == "employee" {
		querySQL += "AND tx_reservation.user_profile_id = " + validation.EscapeString(user_id)
	}

	if resrv_id != "" {
		querySQL += "AND tx_reservation.id = " + validation.EscapeString(resrv_id)
	} else {
		if filter_status != "" {
			querySQL += "AND tx_reservation_status.id = " + validation.EscapeString(filter_status)
		}
		if filter_start_date != "" {
			querySQL += "AND tx_reservation.reservation_date >= " + validation.EscapeString(filter_start_date)
		}
		if filter_end_date != "" {
			querySQL += "AND tx_reservation.reservation_date <= " + validation.EscapeString(filter_end_date)
		}
		if filter_room_id != "" {
			querySQL += "AND mst_room.id = " + validation.EscapeString(filter_room_id)
		}
	}
	querySQL += `
	) AS all_reservations`

	var bjson_reserv []byte
	var rows_reserv_response []response.ReservationResponse

	err := self.db.QueryRow(querySQL).Scan(&bjson_reserv)

	if err != nil {
		return nil, err
	}

	if len(bjson_reserv) == 0 {
		return rows_reserv_response, nil
	}

	err = json.Unmarshal(bjson_reserv, &rows_reserv_response)

	if err != nil {
		return nil, errors.New(fmt.Sprint("Error unmarshaling JSON:", err))
	}

	return rows_reserv_response, nil
}

func NewReservationRepository(db *sql.DB) ReservationRepository {
	return &reservationRepository{
		db: db,
	}
}
