package request

type RoomRequest struct{
	Name string `json:"name" binding:"required"`
	RoomType string `json:"roomtype" binding:"required"`
	Capacity int `json:"capacity" binding:"required"`
}