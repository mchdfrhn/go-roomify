package request

type UserProfileRequest struct {
	Id          string `json:"id"`
	FullName    string `json:"full_name" binding:"required"`
	DivisionId  string `json:"division_id" binding:"required"`
	Address     string `json:"address" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	UserId      string `json:"user_id" binding:"required"`
	RoleId      string `json:"role_id" binding:"required"`
}
