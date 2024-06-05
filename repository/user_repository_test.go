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

var UserProfile1 = model.UserProfile{
	Id:       "550e8400-e29b-41d4-a716-446655440000",
	FullName: "John Doe",
	Division: model.Division{
		Id:   "d123",
		Name: "Engineering",
	},
	Address:     "123 Main St, Springfield",
	PhoneNumber: "555-1234",
	User: model.UserCredential{
		Id:       "u123",
		Username: "johndoe",
		Password: "securepassword",
		Token:    "token123",
	},
	Role: model.Role{
		Id:       "r123",
		Position: "Developer",
	},
}

// Membuat objek UserProfile kedua
var UserProfile2 = model.UserProfile{
	Id:       "650e8400-e29b-41d4-a716-446655440001",
	FullName: "Jane Smith",
	Division: model.Division{
		Id:   "d124",
		Name: "Marketing",
	},
	Address:     "456 Elm St, Springfield",
	PhoneNumber: "555-5678",
	User: model.UserCredential{
		Id:       "u124",
		Username: "janesmith",
		Password: "anothersecurepassword",
		Token:    "token124",
	},
	Role: model.Role{
		Id:       "r124",
		Position: "Manager",
	},
}

// Membuat objek UserProfile ketiga
var UserProfile3 = model.UserProfile{
	Id:       "750e8400-e29b-41d4-a716-446655440002",
	FullName: "Alice Johnson",
	Division: model.Division{
		Id:   "d125",
		Name: "Sales",
	},
	Address:     "789 Oak St, Springfield",
	PhoneNumber: "555-8765",
	User: model.UserCredential{
		Id:       "u125",
		Username: "alicejohnson",
		Password: "yetanothersecurepassword",
		Token:    "token125",
	},
	Role: model.Role{
		Id:       "r125",
		Position: "Salesperson",
	},
}

var UserProfileRequest = request.UserProfileRequest{
	Id:          "750e3400-e29b-41d4-a716-4466g54sv002",
	FullName:    "Mochammad Farhan Ali",
	DivisionId:  "u185",
	Address:     "Sukasarana, Cianjur",
	PhoneNumber: "085771826637",
	UserId:      "u645",
	RoleId:      "r825",
}

// Membuat slice dari UserProfile
var UserProfiles = []model.UserProfile{UserProfile1, UserProfile2, UserProfile3}

type UserRepoTestSuite struct {
	suite.Suite
	repo    UserProfileRepository
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
}

func (suite *UserRepoTestSuite) SetupTest() {
	mockDb, mockSql, _ := sqlmock.New()
	userRepo := NewUserProfileRepository(mockDb)
	suite.mockDb = mockDb
	suite.mockSql = mockSql
	suite.repo = userRepo
}

func (suite *UserRepoTestSuite) TestUserGetList_Success() {
	// Set Column
	rows := suite.mockSql.NewRows([]string{
		"up.id",
		"up.full_name",
		"d.id",
		"d.name",
		"up.address",
		"up.phone_number",
		"u.id",
		"u.username",
		"u.password",
		"u.token",
		"r.id",
		"r.position",
	})

	// Isi Column
	for _, user := range UserProfiles {
		rows.AddRow(
			user.Id,
			user.FullName,
			user.Division.Id,
			user.Division.Name,
			user.Address,
			user.PhoneNumber,
			user.User.Id,
			user.User.Username,
			user.User.Password,
			user.User.Token,
			user.Role.Id,
			user.Role.Position,
		)
	}
	suite.mockSql.ExpectQuery("SELECT up.id,up.full_name,d.id,d.name,up.address,up.phone_number,u.id,u.username,u.password,u.token,r.id,r.position FROM mst_user_profile AS up JOIN mst_division AS d ON up.division_id=d.id JOIN users AS u ON up.user_id=u.id JOIN mst_role AS r ON up.role_id=r.id LIMIT \\$1 OFFSET \\$2").WithArgs(10, 0).WillReturnRows(rows)

	suite.mockSql.ExpectQuery("SELECT COUNT\\(id\\) FROM mst_user_profile").WillReturnRows(sqlmock.NewRows([]string{"COUNT(id)"}).AddRow(3))

	actualUsers, _, err := suite.repo.GetList(1, 10)

	suite.Assert().NoError(err)
	suite.Assert().Equal(3, len(actualUsers))
	suite.Assert().Equal(UserProfile1, actualUsers[0])
}

func (suite *UserRepoTestSuite) TestUserGetList_Fail_Query() {
	expectErr := errors.New("Failed")

	suite.mockSql.ExpectQuery("SELECT up.id,up.full_name,d.id,d.name,up.address,up.phone_number,u.id,u.username,u.password,u.token,r.id,r.position FROM mst_user_profile AS up JOIN mst_division AS d ON up.division_id=d.id JOIN users AS u ON up.user_id=u.id JOIN mst_role AS r ON up.role_id=r.id LIMIT \\$1 OFFSET \\$2").WithArgs(10, 0).WillReturnError(expectErr)

	suite.mockSql.ExpectQuery("SELECT COUNT\\(id\\) FROM mst_user_profile").WillReturnError(expectErr)

	actualUsers, _, err := suite.repo.GetList(1, 10)

	suite.Assert().Equal(expectErr, err)
	suite.Assert().Equal(actualUsers, []model.UserProfile(nil))
}

func (suite *UserRepoTestSuite) TestUserGetList_Fail_Scan_Rows() {
	// Set Column
	expectErr := errors.New("Failed")
	rowsNil := suite.mockSql.NewRows([]string{
		"up.id",
	})

	// Isi Column
	for _, user := range UserProfiles {
		rowsNil.AddRow(
			user.Id,
		)
	}

	suite.mockSql.ExpectQuery("SELECT up.id,up.full_name,d.id,d.name,up.address,up.phone_number,u.id,u.username,u.password,u.token,r.id,r.position FROM mst_user_profile AS up JOIN mst_division AS d ON up.division_id=d.id JOIN users AS u ON up.user_id=u.id JOIN mst_role AS r ON up.role_id=r.id LIMIT \\$1 OFFSET \\$2").WithArgs(10, 0).WillReturnRows(rowsNil)

	actualUsers, _, err := suite.repo.GetList(1, 10)

	suite.Assert().Error(err, expectErr)
	suite.Assert().Equal(actualUsers, []model.UserProfile(nil))
}

func (suite *UserRepoTestSuite) TestUserGetList_Fail_Scan_Count() {
	expectErr := errors.New("Failed")
	rows := suite.mockSql.NewRows([]string{
		"up.id",
		"up.full_name",
		"d.id",
		"d.name",
		"up.address",
		"up.phone_number",
		"u.id",
		"u.username",
		"u.password",
		"u.token",
		"r.id",
		"r.position",
	})

	row := suite.mockSql.NewRows([]string{
		"up.id",
		"up.full_name",
	})

	// Isi Column
	for _, user := range UserProfiles {
		rows.AddRow(
			user.Id,
			user.FullName,
			user.Division.Id,
			user.Division.Name,
			user.Address,
			user.PhoneNumber,
			user.User.Id,
			user.User.Username,
			user.User.Password,
			user.User.Token,
			user.Role.Id,
			user.Role.Position,
		)
	}

	suite.mockSql.ExpectQuery("SELECT up.id,up.full_name,d.id,d.name,up.address,up.phone_number,u.id,u.username,u.password,u.token,r.id,r.position FROM mst_user_profile AS up JOIN mst_division AS d ON up.division_id=d.id JOIN users AS u ON up.user_id=u.id JOIN mst_role AS r ON up.role_id=r.id LIMIT \\$1 OFFSET \\$2").WithArgs(10, 0).WillReturnRows(rows)

	suite.mockSql.ExpectQuery("SELECT COUNT\\(id\\) FROM mst_user_profile").WillReturnRows(row)

	actualUsers, _, err := suite.repo.GetList(1, 10)

	suite.Assert().Error(err, expectErr)
	suite.Assert().Equal(actualUsers, []model.UserProfile(nil))
}

func (suite *UserRepoTestSuite) TestUserGetById_Success() {
	// Set Column
	rows := suite.mockSql.NewRows([]string{
		"up.id",
		"up.full_name",
		"d.id",
		"d.name",
		"up.address",
		"up.phone_number",
		"u.id",
		"u.username",
		"u.password",
		"u.token",
		"r.id",
		"r.position",
	}).AddRow(
		"550e8400-e29b-41d4-a716-446655440000",
		"John Doe",
		"d123",
		"Engineering",
		"123 Main St, Springfield",
		"555-1234",
		"u123",
		"johndoe",
		"securepassword",
		"token123",
		"r123",
		"Developer",
	)

	suite.mockSql.ExpectQuery(`SELECT up\.id,up\.full_name,d\.id,d\.name,up\.address,up\.phone_number,u\.id,u\.username,u\.password,u\.token,r\.id,r\.position FROM mst_user_profile AS up JOIN mst_division AS d ON up\.division_id=d\.id JOIN users AS u ON up\.user_id=u\.id JOIN mst_role AS r ON up\.role_id=r\.id WHERE up\.id = \$1`).WithArgs("550e8400-e29b-41d4-a716-446655440000").WillReturnRows(rows)

	actualUser, err := suite.repo.GetById(UserProfile1.Id)

	suite.Assert().NoError(err)
	suite.Assert().Equal(UserProfile1, actualUser)
}

func (suite *UserRepoTestSuite) TestUserGetById_Fail() {
	expectErr := errors.New("Failed")

	suite.mockSql.ExpectQuery(`SELECT up\.id,up\.full_name,d\.id,d\.name,up\.address,up\.phone_number,u\.id,u\.username,u\.password,u\.token,r\.id,r\.position FROM mst_user_profile AS up JOIN mst_division AS d ON up\.division_id=d\.id JOIN users AS u ON up\.user_id=u\.id JOIN mst_role AS r ON up\.role_id=r\.id WHERE up\.id = \$1`).WithArgs("550e8400-e29b-41d4-a716-446655440000").WillReturnError(expectErr)

	actualUser, err := suite.repo.GetById(UserProfile1.Id)

	suite.Assert().Equal(expectErr, err)
	suite.Assert().Equal(actualUser, model.UserProfile{})
}

func (suite *UserRepoTestSuite) TestUserGetByUsername_Success() {
	// Set Column
	rows := suite.mockSql.NewRows([]string{
		"up.id",
		"up.full_name",
		"d.id",
		"d.name",
		"up.address",
		"up.phone_number",
		"u.id",
		"u.username",
		"u.password",
		"u.token",
		"r.id",
		"r.position",
	}).AddRow(
		"550e8400-e29b-41d4-a716-446655440000",
		"John Doe",
		"d123",
		"Engineering",
		"123 Main St, Springfield",
		"555-1234",
		"u123",
		"johndoe",
		"securepassword",
		"token123",
		"r123",
		"Developer",
	)

	suite.mockSql.ExpectQuery(`SELECT up\.id,up\.full_name,d\.id,d\.name,up\.address,up\.phone_number,u\.id,u\.username,u\.password,u\.token,r\.id,r\.position FROM mst_user_profile AS up JOIN mst_division AS d ON up\.division_id=d\.id JOIN users AS u ON up\.user_id=u\.id JOIN mst_role AS r ON up\.role_id=r\.id WHERE u\.username = \$1`).WithArgs("johndoe").WillReturnRows(rows)

	actualUser, err := suite.repo.GetByUsername(UserProfile1.User.Username)

	suite.Assert().NoError(err)
	suite.Assert().Equal(UserProfile1, actualUser)
}

func (suite *UserRepoTestSuite) TestUserGetByUsername_Fail() {
	expectErr := errors.New("Failed")

	suite.mockSql.ExpectQuery(`SELECT up\.id,up\.full_name,d\.id,d\.name,up\.address,up\.phone_number,u\.id,u\.username,u\.password,u\.token,r\.id,r\.position FROM mst_user_profile AS up JOIN mst_division AS d ON up\.division_id=d\.id JOIN users AS u ON up\.user_id=u\.id JOIN mst_role AS r ON up\.role_id=r\.id WHERE u\.username = \$1`).WithArgs("johndoe").WillReturnError(expectErr)

	actualUser, err := suite.repo.GetByUsername(UserProfile1.User.Username)

	suite.Assert().Equal(expectErr, err)
	suite.Assert().Equal(actualUser, model.UserProfile{})
}

func (suite *UserRepoTestSuite) TestUserCreate_Success() {
	suite.mockSql.ExpectExec("INSERT INTO mst_user_profile").WithArgs(
		"750e3400-e29b-41d4-a716-4466g54sv002",
		"Mochammad Farhan Ali",
		"u185",
		"Sukasarana, Cianjur",
		"085771826637",
		"u645",
		"r825",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	_, err := suite.repo.Create(UserProfileRequest)
	suite.Assert().NoError(err)
}

func (suite *UserRepoTestSuite) TestUserCreate_Fail() {
	expecErr := errors.New("Failed")
	suite.mockSql.ExpectExec("INSERT INTO mst_user_profile").WithArgs(
		"750e3400-e29b-41d4-a716-4466g54sv002",
		"Mochammad Farhan Ali",
		"u185",
		"Sukasarana, Cianjur",
		"085771826637",
		"u645",
		"r825",
	).WillReturnError(expecErr)
	_, err := suite.repo.Create(UserProfileRequest)
	suite.Assert().Equal(expecErr, err)
}

func (suite *UserRepoTestSuite) TestUserUpdate_Success() {
	query := "UPDATE mst_user_profile SET full_name = \\$1,division_id = \\$2,address = \\$3,phone_number = \\$4,role_id = \\$5 WHERE id = \\$6"
	suite.mockSql.ExpectExec(query).
		WithArgs(
			"Mochammad Farhan Ali",
			"u185",
			"Sukasarana, Cianjur",
			"085771826637",
			"r825",
			"750e3400-e29b-41d4-a716-4466g54sv002",
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	_, err := suite.repo.Update(UserProfileRequest)
	suite.Assert().NoError(err)
}

func (suite *UserRepoTestSuite) TestUserUpdate_Fail() {
	expectErr := errors.New("Failed")
	query := "UPDATE mst_user_profile SET full_name = \\$1,division_id = \\$2,address = \\$3,phone_number = \\$4,role_id = \\$5 WHERE id = \\$6"
	suite.mockSql.ExpectExec(query).
		WithArgs(
			"Mochammad Farhan Ali",                 
			"u185",                                 
			"Sukasarana, Cianjur",                  
			"085771826637",                         
			"r825",                                 
			"750e3400-e29b-41d4-a716-4466g54sv002", 
		).
		WillReturnError(expectErr)
	_, err := suite.repo.Update(UserProfileRequest)
	suite.Assert().Equal(expectErr, err)
}

func (suite *UserRepoTestSuite) TestUserDelete_Success() {
	query := "DELETE FROM mst_user_profile WHERE id = \\$1"
	suite.mockSql.ExpectExec(query).
		WithArgs("550e8400-e29b-41d4-a716-446655440000").WillReturnResult(sqlmock.NewResult(1, 1))
	err := suite.repo.Delete(UserProfile1.Id)
	suite.Assert().NoError(err)
}

func (suite *UserRepoTestSuite) TestUserDelete_Fail() {
	expectErr := errors.New("Failed")
	query := "DELETE FROM mst_user_profile WHERE id = \\$1"
	suite.mockSql.ExpectExec(query).
		WithArgs("550e8400-e29b-41d4-a716-446655440000").WillReturnError(expectErr)
	err := suite.repo.Delete(UserProfile1.Id)
	suite.Assert().Equal(expectErr, err)
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}
