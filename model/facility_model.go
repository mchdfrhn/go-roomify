package model

type Facility struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	IsAvailable   bool   `json:"is_available"`
	IsReserveable bool   `json:"is_reserveable"`
	RoomId        string `json:"room_id"`
}
