package request

type RoomRequest struct {
	Name          string `json:"name" binding:"required"`
	RoomTypeId    string `json:"room_type_id" binding:"required"`
	Capacity      int    `json:"capacity" binding:"required"`
	IsReserveable *bool  `json:"is_reserveable" binding:"required"`
}

type RoomStatusRequest struct {
	Id          string `json:"id" binding:"required"`
	IsAvailable *bool  `json:"is_available" binding:"required"`
}
