package model

type UserCredential struct {
	Id			string	`json:"id"`
	Username 	string	`json:"username"`
	Password 	string	`json:"password"`
	Role 		string	`json:"role"`
	Token 		string	`json:"token"`
}
