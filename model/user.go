package model

type UserProfile struct {
	Id          string         `json:"id" binding:"required"`
	FullName    string         `json:"full_name" binding:"required"`
	Division    Division       `json:"division" binding:"required"`
	Address     string         `json:"address" binding:"required"`
	PhoneNumber string         `json:"phone_number" binding:"required"`
	User        UserCredential `json:"user_credential" binding:"required"`
	Role        Role           `json:"role" binding:"required"`
}
