package repository

import (
	"database/sql"
	"errors"
	"go-roomify/model"
	"go-roomify/model/dto/request"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

var user1 = model.UserCredential{
	Id:       "4f4334e2-0465-4fd6-bfb5-91683b574945",
	Username: "user1",
	Password: "password1",
	Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoidXNlcjEiLCJleHAiOjE3MTc1MDA1Mzh9.kyBRm7v9FmWFnHDiDAXDCgOJWGuhvLY4G5MWkJYiKs8",
}

var user2 = model.UserCredential{
	Id:       "5ffcbe6f-6037-4e32-933d-1a69001a46ec",
	Username: "user2",
	Password: "password2",
	Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoidXNlcjIiLCJleHAiOjE3MTc1MDA1Mzh9.3BGFhKeOQ25zVZMfwktjDALf4s19ZjIEgqu8z0jSj38",
}

var user3 = model.UserCredential{
	Id:       "d0b886b2-0389-4fd5-b2c2-63af5abdd1fb",
	Username: "user3",
	Password: "password3",
	Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoidXNlcjMiLCJleHAiOjE3MTc1MDA1Mzh9.hqYdPX4NAajLthD910cc99FLGYTX1qAjb9p-g0ba4oA",
}

var userJwt = model.UserCredentialJwt{
	Id:       "4f4334e2-0465-4fd6-bfb5-91683b574945",
	Username: "user1",
	Password: "password1",
	Role:     "admin",
	Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoidXNlcjEiLCJleHAiOjE3MTc1MDA1Mzh9.kyBRm7v9FmWFnHDiDAXDCgOJWGuhvLY4G5MWkJYiKs8",
}

var userReq = request.UserUpdatePasswordRequest{
	Id:          "4f4334e2-0465-4fd6-bfb5-91683b574945",
	OldPassword: "oldPassword",
	NewPassword: "newPassword",
}

var users = []model.UserCredential{user1, user2, user3}

type UserCredentialRepoTestSuite struct {
	suite.Suite
	repo    UserCredentialRepository
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
}

func (suite *UserCredentialRepoTestSuite) TestUserCredentialGetByUsername_Success() {
	// Set Column
	rows := suite.mockSql.NewRows([]string{
		"users.id",
		"users.username",
		"users.password",
		"mst_role.position",
		"users.token",
	}).AddRow(
		"4f4334e2-0465-4fd6-bfb5-91683b574945",
		"user1",
		"password1",
		"admin",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoidXNlcjEiLCJleHAiOjE3MTc1MDA1Mzh9.kyBRm7v9FmWFnHDiDAXDCgOJWGuhvLY4G5MWkJYiKs8",
	)

	suite.mockSql.ExpectQuery(`SELECT users.id,users.username,users.password,mst_role.position,users.token FROM mst_user_profile JOIN mst_role ON mst_user_profile.role_id = mst_role.id JOIN users ON mst_user_profile.user_id = users.id WHERE users.username = \$1 LIMIT \$2`).WithArgs("user1", 1).WillReturnRows(rows)

	actualUser, err := suite.repo.GetByUsername(userJwt.Username)

	suite.Assert().NoError(err)
	suite.Assert().Equal(userJwt, actualUser)
}

func (suite *UserCredentialRepoTestSuite) TestUserCredentialGetByUsername_Fail() {
	expectErr := errors.New("Failed")

	suite.mockSql.ExpectQuery(`SELECT users.id,users.username,users.password,mst_role.position,users.token FROM mst_user_profile JOIN mst_role ON mst_user_profile.role_id = mst_role.id JOIN users ON mst_user_profile.user_id = users.id WHERE users.username = \$1 LIMIT \$2`).WithArgs("user1", 1).WillReturnError(expectErr)

	actualUser, err := suite.repo.GetByUsername(userJwt.Username)

	suite.Assert().Equal(expectErr, err)
	suite.Assert().Equal(actualUser, model.UserCredentialJwt{})
}

func (suite *UserCredentialRepoTestSuite) TestUserCredentialGetList_Success() {
	// Set Column
	rows := suite.mockSql.NewRows([]string{
		"id",
		"username",
		"password",
		"token",
	})

	// Isi Column
	for _, user := range users {
		rows.AddRow(
			user.Id,
			user.Username,
			user.Password,
			user.Token,
		)
	}
	suite.mockSql.ExpectQuery(`SELECT id,username,password,token FROM users LIMIT \$1 OFFSET \$2`).WithArgs(10, 0).WillReturnRows(rows)

	suite.mockSql.ExpectQuery(`SELECT COUNT\(id\) FROM users`).WillReturnRows(sqlmock.NewRows([]string{"COUNT(id)"}).AddRow(3))

	actualUsers, _, err := suite.repo.GetList(1, 10)

	suite.Assert().NoError(err)
	suite.Assert().Equal(3, len(actualUsers))
	suite.Assert().Equal(user1, actualUsers[0])
}

func (suite *UserCredentialRepoTestSuite) TestUserCredentialGetList_Fail_Query() {
	expectErr := errors.New("Failed")

	suite.mockSql.ExpectQuery(`SELECT id,username,password,token FROM users LIMIT \$1 OFFSET \$2`).WithArgs(10, 0).WillReturnError(expectErr)

	suite.mockSql.ExpectQuery(`SELECT COUNT\(id\) FROM users`).WillReturnError(expectErr)

	actualUsers, _, err := suite.repo.GetList(1, 10)

	suite.Assert().Equal(expectErr, err)
	suite.Assert().Equal(actualUsers, []model.UserCredential(nil))
}

func (suite *UserCredentialRepoTestSuite) TestUserCredentialGetList_Fail_Scan_Rows() {
	// Set Column
	expectErr := errors.New("Failed")
	rowsNil := suite.mockSql.NewRows([]string{
		"id",
	})

	// Isi Column
	for _, user := range users {
		rowsNil.AddRow(
			user.Id,
		)
	}

	suite.mockSql.ExpectQuery(`SELECT id,username,password,token FROM users LIMIT \$1 OFFSET \$2`).WithArgs(10, 0).WillReturnRows(rowsNil)

	actualUsers, _, err := suite.repo.GetList(1, 10)

	suite.Assert().Error(err, expectErr)
	suite.Assert().Equal(actualUsers, []model.UserCredential(nil))
}

func (suite *UserCredentialRepoTestSuite) TestUserCredentialGetById_Success() {
	// Set Column
	rows := suite.mockSql.NewRows([]string{
		"users.id",
		"users.username",
		"users.password",
		"users.token",
	}).AddRow(
		"4f4334e2-0465-4fd6-bfb5-91683b574945",
		"user1",
		"password1",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoidXNlcjEiLCJleHAiOjE3MTc1MDA1Mzh9.kyBRm7v9FmWFnHDiDAXDCgOJWGuhvLY4G5MWkJYiKs8",
	)

	suite.mockSql.ExpectQuery(`SELECT id,username,password,token FROM users WHERE id = \$1`).WithArgs("4f4334e2-0465-4fd6-bfb5-91683b574945").WillReturnRows(rows)

	actualUser, err := suite.repo.GetById(user1.Id)

	suite.Assert().NoError(err)
	suite.Assert().Equal(user1, actualUser)
}

func (suite *UserCredentialRepoTestSuite) TestUserCredentialGetById_Fail() {
	expectErr := errors.New("Failed")

	suite.mockSql.ExpectQuery(`SELECT id,username,password,token FROM users WHERE id = \$1`).WithArgs("4f4334e2-0465-4fd6-bfb5-91683b574945").WillReturnError(expectErr)

	actualUser, err := suite.repo.GetById(user1.Id)

	suite.Assert().Equal(expectErr, err)
	suite.Assert().Equal(actualUser, model.UserCredential{})
}

func (suite *UserCredentialRepoTestSuite) TestUserCredentialCreate_Success() {
	suite.mockSql.ExpectExec("INSERT INTO users").WithArgs(
		"4f4334e2-0465-4fd6-bfb5-91683b574945",
		"user1",
		"password1",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoidXNlcjEiLCJleHAiOjE3MTc1MDA1Mzh9.kyBRm7v9FmWFnHDiDAXDCgOJWGuhvLY4G5MWkJYiKs8",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	_, err := suite.repo.AddNew(user1)
	suite.Assert().NoError(err)
}

func (suite *UserCredentialRepoTestSuite) TestUserCredentialCreate_Fail() {
	expecErr := errors.New("Failed")
	suite.mockSql.ExpectExec("INSERT INTO users").WithArgs(
		"4f4334e2-0465-4fd6-bfb5-91683b574945",
		"user1",
		"password1",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoidXNlcjEiLCJleHAiOjE3MTc1MDA1Mzh9.kyBRm7v9FmWFnHDiDAXDCgOJWGuhvLY4G5MWkJYiKs8",
	).WillReturnError(expecErr)
	_, err := suite.repo.AddNew(user1)
	suite.Assert().Equal(expecErr, err)
}

func (suite *UserCredentialRepoTestSuite) TestUserCredentialUpdate_Success() {
	suite.mockSql.ExpectExec(`UPDATE users SET password = \$1 WHERE id = \$2`).
		WithArgs(
			"newPassword",
			"4f4334e2-0465-4fd6-bfb5-91683b574945",
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err := suite.repo.UpdatePassword(userReq)
	suite.Assert().NoError(err)
}

func (suite *UserCredentialRepoTestSuite) TestUserCredentialUpdate_Fail() {
	expectErr := errors.New("Failed")
	suite.mockSql.ExpectExec(`UPDATE users SET password = \$1 WHERE id = \$2`).
		WithArgs(
			"newPassword",
			"4f4334e2-0465-4fd6-bfb5-91683b574945",
		).
		WillReturnError(expectErr)
	err := suite.repo.UpdatePassword(userReq)
	suite.Assert().Equal(expectErr, err)
}

func (suite *UserCredentialRepoTestSuite) TestUserCredentialDelete_Success() {
	suite.mockSql.ExpectExec(`DELETE FROM users WHERE id = \$1`).
		WithArgs("4f4334e2-0465-4fd6-bfb5-91683b574945").WillReturnResult(sqlmock.NewResult(1, 1))
	err := suite.repo.Delete(user1.Id)
	suite.Assert().NoError(err)
}

func (suite *UserCredentialRepoTestSuite) TestUserCredentialDelete_Fail() {
	expectErr := errors.New("Failed")
	suite.mockSql.ExpectExec(`DELETE FROM users WHERE id = \$1`).
		WithArgs("4f4334e2-0465-4fd6-bfb5-91683b574945").WillReturnError(expectErr)
	err := suite.repo.Delete(user1.Id)
	suite.Assert().Equal(expectErr, err)
}

func (suite *UserCredentialRepoTestSuite) SetupTest() {
	mockDb, mockSql, _ := sqlmock.New()
	userCredentialRepo := NewUserCredentialRepository(mockDb)
	suite.mockDb = mockDb
	suite.mockSql = mockSql
	suite.repo = userCredentialRepo
}

func TestUserCredentialRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserCredentialRepoTestSuite))
}