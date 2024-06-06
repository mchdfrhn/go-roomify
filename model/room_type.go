package model

type RoomType struct{
	Id string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}