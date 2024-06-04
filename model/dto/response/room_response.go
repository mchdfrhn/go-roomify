package response

import "go-roomify/model"

type RoomResponse struct {
	Room       model.Room
	Facilities []model.Facility
}
