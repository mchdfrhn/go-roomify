package model

type Room struct{
	Id string `json:"id"`
	Name string `json:"name"`
	RoomType string `json:"roomtype"`
	Capacity int `json:"capacity"`
	IsAvailable bool `json:"is_available"`
}