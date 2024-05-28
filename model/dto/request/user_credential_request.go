package request 

type UserCredentialRequest struct {
	Username string	`json:"username" binding:"required"`
	Password string	`json:"password" binding:"required"`
}

type UserUpdatePasswordRequest struct {
	Id string	`json:"id" binding:"required"`
	OldPassword string	`json:"old_password" binding:"required"`
	NewPassword string	`json:"new_password" binding:"required"`
}


