package model

type Room struct{
	Id string `json:"id"`
	Name string `json:"name"`
	RoomType string `json:"roomtype"`
	Capacity string `json:"capacity"`
	IsAvailable string `json:"is_available"`
}