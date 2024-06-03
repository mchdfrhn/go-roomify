package request

type ReservationStatusRequest struct {
	ReservationId string `json:"reservation_id" binding:"required"`
	StatusId      string `json:"status_id"`
	Description   string `json:"description"`
}

type ReservationGetListFilter struct {
	UserId          string
	UserRole        string
	ReservationId   string
	FilterStatus    string
	FilterStartDate string
	FilterEndDate   string
}
