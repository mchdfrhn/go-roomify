package model

type Room struct {
	Id            string `json:"id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	RoomTypeId    string `json:"room_type_id" binding:"required"`
	Capacity      int    `json:"capacity" binding:"required"`
	IsAvailable   *bool  `json:"is_available" binding:"required"`
	IsReserveable *bool  `json:"is_reserveable" binding:"required"`
}
