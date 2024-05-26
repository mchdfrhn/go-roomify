package repository

import (
	"database/sql"
	"go-roomify/model"
	"go-roomify/utils/query"
	"fmt"
)


type UserCredentialRepository interface {
	GetByUsername(username string) (model.UserCredential, error)
}

type userCredentialRepository struct {
	db *sql.DB
}

func (self *userRepository) GetByUsername(username string) (model.UserCredential, error) {
	var user_cr model.UserCredential

	qselect := query.QSelect{DB: self.db}

	qselect.Table("mst_user_profile")
	qselect.Column(
		"users.id",
		"users.username",
		"users.password",
		"users.token")
	qselect.Join("users", "mst_user_profile.users_id = users.id")
	qselect.Where("users.username", "=", username)
	qselect.Limit(1)

	rows, err := qselect.Run()

	for rows.Next() {
		
	}

	if err != nil {
		return model.UserCredential{},err
	}

	return user,nil
}
