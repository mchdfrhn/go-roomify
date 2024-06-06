package request

type RoomRequest struct {
	Name     string `json:"name" binding:"required"`
	RoomType string `json:"roomtype" binding:"required"`
	Capacity int    `json:"capacity" binding:"required"`
	IsReserveable *bool `json:"is_reserveable" binding:"required"`
}

type RoomStatusRequest struct {
	Id          string `json:"id" binding:"required"`
	IsAvailable *bool  `json:"is_available" binding:"required"`
}
