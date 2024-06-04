package usecase

import (
	"errors"
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/model/dto/request"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

type RepoMock struct {
	mock.Mock
}

func (r *RepoMock) GetList(page, size int) ([]model.UserProfile, dto.Paging, error) {
	args := r.Called(page, size)
	if args.Get(2) != nil {
		return []model.UserProfile{}, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.UserProfile), dto.Paging{}, nil
}

func (r *RepoMock) GetById(id string) (model.UserProfile, error) {
	args := r.Called(id)
	if args.Get(1) != nil {
		return model.UserProfile{}, args.Error(1)
	}
	return args.Get(0).(model.UserProfile), nil
}

func (r *RepoMock) GetByUsername(username string) (model.UserProfile, error) {
	args := r.Called(username)
	if args.Get(1) != nil {
		return model.UserProfile{}, args.Error(1)
	}
	return args.Get(0).(model.UserProfile), nil
}

func (r *RepoMock) Create(newUser request.UserProfileRequest) (request.UserProfileRequest, error) {
	args := r.Called(newUser)
	if args.Get(1) != nil {
		return request.UserProfileRequest{}, args.Error(1)
	}
	return args.Get(0).(request.UserProfileRequest), nil
}

func (r *RepoMock) Update(newUser request.UserProfileRequest) (request.UserProfileRequest, error) {
	args := r.Called(newUser)
	if args.Get(1) != nil {
		return request.UserProfileRequest{}, args.Error(1)
	}
	return args.Get(0).(request.UserProfileRequest), nil
}

func (r *RepoMock) Delete(id string) error {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

type UserrUsecaseTestSuite struct {
	suite.Suite
	repoMock *RepoMock
	usecase  UserProfileUsecase
}

func (suite *UserrUsecaseTestSuite) TestUserGetList_Success() {
	suite.repoMock.On("GetList", 10, 0).Return(UserProfiles, dto.Paging{}, nil)
	actual, _, err := suite.usecase.GetList(10, 0)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), UserProfiles, actual)
}

func (suite *UserrUsecaseTestSuite) TestUserGetList_Fail() {
	expectErr := errors.New("Failed")
	suite.repoMock.On("GetList", 10, 0).Return([]model.UserProfile{}, dto.Paging{}, expectErr)
	_, _, err := suite.usecase.GetList(10, 0)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectErr, err)
}

func (suite *UserrUsecaseTestSuite) TestUserGetById_Success() {
	suite.repoMock.On("GetById", UserProfile1.Id).Return(UserProfile1, nil)
	actual, err := suite.usecase.GetById(UserProfile1.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), UserProfile1, actual)
}

func (suite *UserrUsecaseTestSuite) TestUserGetById_Fail() {
	expectErr := errors.New("Failed")
	suite.repoMock.On("GetById", UserProfile1.Id).Return(model.UserProfile{}, expectErr)
	_, err := suite.usecase.GetById(UserProfile1.Id)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectErr, err)
}

func (suite *UserrUsecaseTestSuite) TestUserGetByUsername_Success() {
	suite.repoMock.On("GetByUsername", UserProfile1.User.Username).Return(UserProfile1, nil)
	actual, err := suite.usecase.GetByUsername(UserProfile1.User.Username)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), UserProfile1, actual)
}

func (suite *UserrUsecaseTestSuite) TestUserGetByUsername_Fail() {
	expectErr := errors.New("Failed")
	suite.repoMock.On("GetByUsername", UserProfile1.User.Username).Return(model.UserProfile{}, expectErr)
	_, err := suite.usecase.GetByUsername(UserProfile1.User.Username)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectErr, err)
}

func (suite *UserrUsecaseTestSuite) TestUserCreate_Success() {
	suite.repoMock.On("Create", UserProfileRequest).Return(UserProfileRequest, nil)
	actual, err := suite.usecase.Create(UserProfileRequest)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), UserProfileRequest, actual)
}

func (suite *UserrUsecaseTestSuite) TestUserCreate_Error() {
	expectErr := errors.New("Failed")
	suite.repoMock.On("Create", UserProfileRequest).Return(request.UserProfileRequest{}, expectErr)
	_, err := suite.usecase.Create(UserProfileRequest)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectErr, err)
}

func (suite *UserrUsecaseTestSuite) TestUserUpdate_Success() {
	suite.repoMock.On("Update", UserProfileRequest).Return(UserProfileRequest, nil)
	actual, err := suite.usecase.Update(UserProfileRequest)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), UserProfileRequest, actual)
}

func (suite *UserrUsecaseTestSuite) TestUserUpdate_Error() {
	expectErr := errors.New("Failed")
	suite.repoMock.On("Update", UserProfileRequest).Return(request.UserProfileRequest{}, expectErr)
	_, err := suite.usecase.Update(UserProfileRequest)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectErr, err)
}

func (suite *UserrUsecaseTestSuite) TestUserDelete_Success() {
	suite.repoMock.On("Delete", UserProfileRequest.Id).Return(nil)
	err := suite.usecase.Delete(UserProfileRequest.Id)
	assert.Nil(suite.T(), err)
}

func (suite *UserrUsecaseTestSuite) TestUserDelete_Error() {
	expectErr := errors.New("Failed")
	suite.repoMock.On("Delete", UserProfileRequest.Id).Return(expectErr)
	err := suite.usecase.Delete(UserProfileRequest.Id)
	assert.Error(suite.T(), err)
}

func (suite *UserrUsecaseTestSuite) SetupTest() {
	suite.repoMock = new(RepoMock)
	suite.usecase = NewUserUsecase(suite.repoMock)
}

func TestUserrUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserrUsecaseTestSuite))
}
