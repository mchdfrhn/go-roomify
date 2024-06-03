package request

type ReservationStatusRequest struct {
	ReservationId string `json:"reservation_id" binding:"required"`
	StatusId      string `json:"status_id"`
	Description   string `json:"description"`
}

type ReservationGetListFilter struct {
	UserId       string
	UserRole     string
	FilterStatus string
}
