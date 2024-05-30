package request

type ReservationStatus struct {
	Id          string `json:"id" binding:"required"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}
