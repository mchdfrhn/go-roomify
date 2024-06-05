package model

type Facility struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	RoomId        string `json:"room_id"`
	IsAvailable   bool   `json:"is_available"`
	IsReserveable bool   `json:"is_reserveable"`
}
