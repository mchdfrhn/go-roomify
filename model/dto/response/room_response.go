package response

import "go-roomify/model"

type RoomResponse struct {
	Id            string `json:"id"`
    Name          string `json:"name"`
    RoomType      string `json:"roomtype"`
    Capacity      int    `json:"capacity"`
    IsAvailable   bool   `json:"is_available"`
    IsReserveable bool   `json:"is_reserveable"`
	Facilities []model.Facility `json:"facilities"`
}
