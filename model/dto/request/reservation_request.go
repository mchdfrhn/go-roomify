package request

type ReservationRequest struct {
	Id                 string                     `json:"id"`
	UserProfileId      string                     `json:"user_profile_id"`
	ReservationDate    string                     `json:"reservation_date"`
	StartDate          string                     `json:"start_date" binding:"required"`
	EndDate            string                     `json:"end_date" binding:"required"`
	StatusId           string                     `json:"status_id"`
	RoomId             string                     `json:"room_id"  binding:"required"`
	RequestMessage     string                     `json:"request_message"`
	ResponseMessage    string                     `json:"response_message"`
	AdditionalFacility []ReservationDetailRequest `json:"additional_facility"`
}

type ReservationDetailRequest struct {
	Id            string `json:"id"`
	ReservationId string `json:"reservation_id"`
	FacilityId    string `json:"facility_id"`
}

type ReservationStatusRequest struct {
	ReservationId   string `json:"reservation_id" binding:"required"`
	StatusId        string `json:"status_id"`
	ResponseMessage string `json:"response_message"`
}

type ReservationGetListFilter struct {
	UserId           string
	UserRole         string
	ReservationId    string
	FilterStatus     string
	FilterStartDate  string
	FilterEndDate    string
	FilterRoomId     string
	FilterPageNumber int
	FilterPageSize   int
}
