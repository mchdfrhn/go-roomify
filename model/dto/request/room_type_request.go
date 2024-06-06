package request

type RoomTypeRequest struct {
	Name     string `json:"name" binding:"required"`
}
