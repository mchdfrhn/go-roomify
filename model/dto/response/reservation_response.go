package response

import "go-roomify/model"

type ReservationResponse struct {
	ID                 string                         `json:"id"`
	UserProfile        ReservationUserProfileResponse `json:"user_profile"`
	ReservationDate    string                         `json:"reservation_date"`
	StartTime          string                         `json:"start_time"`
	EndTime            string                         `json:"end_time"`
	ReservationStatus  model.ReservationStatus        `json:"reservation_status"`
	Room               ReservationRoomResponse        `json:"room"`
	RequestMessage     string                         `json:"request_message"`
	ResponseMessage    string                         `json:"response_message"`
	AdditionalFacility []model.Facility               `json:"additional_facility"`
}

type ReservationUserProfileResponse struct {
	ID          string         `json:"id"`
	FullName    string         `json:"full_name"`
	Division    model.Division `json:"division"`
	Role        model.Role     `json:"role"`
	Address     string         `json:"address"`
	PhoneNumber string         `json:"phone_number"`
}

type ReservationRoomResponse struct {
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	RoomType      model.RoomType   `json:"roomtype"`
	Capacity      int              `json:"capacity"`
	IsAvailable   bool             `json:"is_available"`
	IsReserveable bool             `json:"is_reserveable"`
	Facility      []model.Facility `json:"facility"`
}

// type ReservationDetailResponse struct {
// 	Id                 string           `json:"id"`
// 	ReservationId      string           `json:"reservation_id"`
// 	AdditionalFacility []model.Facility `json:"additional_facility"`
// }
