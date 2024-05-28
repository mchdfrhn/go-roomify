package response

type UserCredentialResponse struct {
	AccessToken string		`json:"access_token"`
	UserId		string 		`json:"user_id"`
}
