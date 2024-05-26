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

func (self *userCredentialRepository) GetByUsername(username string) (model.UserCredential, error) {
	qselect := query.QSelect{DB: self.db}

	qselect.Table("mst_user_profile")
	qselect.Column(
		"users.id",
		"users.username",
		"users.password",
		"mst_role.position",
		"users.token")
	qselect.Join("mst_role", "mst_user_profile.role_id = mst_role.id")
	qselect.Join("users", "mst_user_profile.user_id = users.id")
	qselect.Where("users.username", "=", username)
	qselect.Limit(1)

	fmt.Println("Query: ", qselect.GetQuery())

	rows, err := qselect.Run()

	if err != nil {
		return model.UserCredential{}, err
	}

	var r_user_cr model.UserCredential

	for rows.Next() {
		err := rows.Scan(
			&r_user_cr.Id,
			&r_user_cr.Username,
			&r_user_cr.Password,
			&r_user_cr.Role,
			&r_user_cr.Token)

		if err != nil {
			return model.UserCredential{}, err
		}

		fmt.Println("Username: ", r_user_cr.Username, r_user_cr.Password)
	}

	rows.Close()

	return r_user_cr, nil
}

func NewUserCredentialRepository(db *sql.DB) (UserCredentialRepository) {
	return &userCredentialRepository {
		db: db,
	}
}