package model

type Room struct {
	Id            string `json:"id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	RoomType      string `json:"roomtype" binding:"required"`
	Capacity      int    `json:"capacity" binding:"required"`
	IsAvailable   *bool  `json:"is_available" binding:"required"`
	IsReserveable *bool  `json:"is_reserveable" binding:"required"`
}
