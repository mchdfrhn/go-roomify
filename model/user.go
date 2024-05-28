package model

type UserProfile struct {
	Id          string   `json:"id" binding:"required"`
	FullName    string   `json:"full_name" binding:"required"`
	Division    Division `json:"division" binding:"required"`
	Address     string   `json:"address" binding:"required"`
	PhoneNumber string   `json:"phone_number" binding:"required"`
	User        User     `json:"user" binding:"required"`
	Role        Role     `json:"role" binding:"required"`
}

// type Division struct {
// 	Id   string `json:"id" binding:"required"`
// 	Name string `json:"name" binding:"required"`
// }

type User struct {
	Id       string `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

// type Role struct {
// 	Id string `json:"id" binding:"required"`
// 	Position string `json:"position" binding:"required"`
// }
