package response

type FacilityForRoomResponse struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	RoomId string `json:"room_id"`
}

type RoomResponse struct {
	Id            string                    `json:"id"`
	Name          string                    `json:"name"`
	RoomType      string                    `json:"roomtype"`
	Capacity      int                       `json:"capacity"`
	IsAvailable   bool                      `json:"is_available"`
	IsReserveable bool                      `json:"is_reserveable"`
	Facilities    []FacilityForRoomResponse `json:"facilities"`
}
