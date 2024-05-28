package response

import (
	"go-roomify/model"
)

type ReservationResponse struct {
	Id					string	`json:"id"`
	UserProfileId 		string	`json:"user_profile_id"`
	ReservationDate 	string	`json:"reservation_date"`
	StartDate		 	string	`json:"start_date"`
	EndDate		 		string	`json:"end_date"`
	Status 				int	`json:"status"`
	Description			string  `json:"description"`
	Detail []ReservationDetailResponse `json:"detail"`
}

type ReservationDetailResponse struct {
	Id					string		`json:"id"`
	ReservationId 		string 		`json:"reservation_id"`
	Room 				model.Room	`json:"room"`
	Equipment			string 		`json:"equipment"`
}