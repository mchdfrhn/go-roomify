package repository_test

// import (
// 	"database/sql"
// 	"errors"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/your_project/model"
// 	"github.com/your_project/repository"
// 	"github.com/your_project/query"
// )

// func TestGetByUsername(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	userRepo := &repository.UserCredentialRepository{DB: db}
// 	username := "test_user"
// 	expectedQuery := `SELECT users.id, users.username, users.password, mst_role.position, users.token FROM mst_user_profile JOIN mst_role ON mst_user_profile.role_id = mst_role.id JOIN users ON mst_user_profile.user_id = users.id WHERE users.username = \\$1 LIMIT \\$2`

// 	// Test case for successful retrieval of user credentials
// 	t.Run("successful retrieval", func(t *testing.T) {
// 		columns := []string{"id", "username", "password", "position", "token"}
// 		mock.ExpectQuery(expectedQuery).WithArgs(username).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, "test_user", "hashed_password", "admin", "jwt_token"))

// 		result, err := userRepo.GetByUsername(username)
// 		assert.NoError(t, err)
// 		assert.Equal(t, model.UserCredentialJwt{
// 			Id:       1,
// 			Username: "test_user",
// 			Password: "hashed_password",
// 			Role:     "admin",
// 			Token:    "jwt_token",
// 		}, result)
// 	})

// 	// Test case for error during query execution
// 	t.Run("query execution error", func(t *testing.T) {
// 		mock.ExpectQuery(expectedQuery).WithArgs(username).WillReturnError(errors.New("query execution error"))

// 		_, err := userRepo.GetByUsername(username)
// 		assert.Error(t, err)
// 		assert.Equal(t, "query execution error", err.Error())
// 	})

// 	// Test case for no rows found
// 	t.Run("no rows found", func(t *testing.T) {
// 		columns := []string{"id", "username", "password", "position", "token"}
// 		mock.ExpectQuery(expectedQuery).WithArgs(username).WillReturnRows(sqlmock.NewRows(columns))

// 		result, err := userRepo.GetByUsername(username)
// 		assert.NoError(t, err)
// 		assert.Equal(t, model.UserCredentialJwt{}, result)
// 	})

// 	// Test case for error during row scanning
// 	t.Run("row scan error", func(t *testing.T) {
// 		columns := []string{"id", "username", "password", "position", "token"}
// 		mock.ExpectQuery(expectedQuery).WithArgs(username).WillReturnRows(sqlmock.NewRows(columns).AddRow("invalid_id", "test_user", "hashed_password", "admin", "jwt_token"))

// 		_, err := userRepo.GetByUsername(username)
// 		assert.Error(t, err)
// 	})

// 	// Ensure all expectations were met
// 	err = mock.ExpectationsWereMet()
// 	assert.NoError(t, err)
// }
