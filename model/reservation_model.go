package model

type Reservation struct {
	Id              string            `json:"id"`
	UserProfileId   string            `json:"user_profile_id"`
	ReservationDate string            `json:"reservation_date"`
	StartDate       string            `json:"start_date" binding:"required"`
	EndDate         string            `json:"end_date" binding:"required"`
	Status          ReservationStatus `json:"status"`
	Description     string            `json:"description"`
	Detail          ReservationDetail `json:"detail" binding:"required"`
}

type ReservationDetail struct {
	Id            string `json:"id"`
	ReservationId string `json:"reservation_id"`
	RoomId        string `json:"room_id" binding:"required"`
	Equipment     string `json:"equipment"`
}

type ReservationStatus struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
