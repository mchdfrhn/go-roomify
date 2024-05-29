package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-roomify/model"
	"go-roomify/model/dto"
	"go-roomify/model/dto/request"
	"go-roomify/model/dto/response"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	User: model.User{
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
	User: model.User{
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
	User: model.User{
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

type UsecaseMock struct {
	mock.Mock
}

type UserControllerTestSuite struct {
	suite.Suite
	routerMock  *gin.Engine
	usecaseMock *UsecaseMock
	controller  *UserController
}

func (u *UsecaseMock) GetList(page, size int) ([]model.UserProfile, dto.Paging, error) {
	return []model.UserProfile{}, dto.Paging{}, nil
}

func (u *UsecaseMock) GetById(id string) (model.UserProfile, error) {
	return model.UserProfile{}, nil
}

func (u *UsecaseMock) GetByUsername(username string) (model.UserProfile, error) {
	return model.UserProfile{}, nil
}

func (u *UsecaseMock) Create(newUser request.UserProfileRequest) (request.UserProfileRequest, error) {
	args := u.Called(newUser)
	if args.Get(1) != nil {
		return request.UserProfileRequest{}, args.Error(1)
	}
	return args.Get(0).(request.UserProfileRequest), nil
}

func (u *UsecaseMock) Update(newUser request.UserProfileRequest) (request.UserProfileRequest, error) {
	args := u.Called(newUser)
	if args.Get(1) != nil {
		return request.UserProfileRequest{}, args.Error(1)
	}
	return args.Get(0).(request.UserProfileRequest), nil
}

func (u *UsecaseMock) Delete(id string) error {
	args := u.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (suite *UserControllerTestSuite) TestRegisterHandler_Success() {
	// Set SEED
	my_rand := rand.New(rand.NewSource(42))
	uuid.SetRand(my_rand)
	UserProfileRequest.Id = uuid.NewString()

	// Set SEED 2 - Untuk digunakan di Customer Controller
	my_rand2 := rand.New(rand.NewSource(42))
	uuid.SetRand(my_rand2)

	suite.usecaseMock.On("Create", UserProfileRequest).Return(UserProfileRequest, nil)
	recorder := httptest.NewRecorder()
	reqBody, _ := json.Marshal(UserProfileRequest)
	request, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(reqBody))
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	suite.controller.registerHandler(ctx)

	var response response.SingleResponse
	json.Unmarshal([]byte(recorder.Body.Bytes()), &response)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.NotEmpty(suite.T(), response.Data)
}

func (suite *UserControllerTestSuite) TestRegisterHandler_FailUsecase() {
	// Set SEED
	my_rand := rand.New(rand.NewSource(42))
	uuid.SetRand(my_rand)
	UserProfileRequest.Id = uuid.NewString()

	// Set SEED 2 - Untuk digunakan di Customer Controller
	my_rand2 := rand.New(rand.NewSource(42))
	uuid.SetRand(my_rand2)

	expectErr := errors.New("Failed")
	suite.usecaseMock.On("Create", UserProfileRequest).Return(request.UserProfileRequest{}, expectErr)
	recorder := httptest.NewRecorder()
	reqBody, _ := json.Marshal(UserProfileRequest)
	request, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(reqBody))
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	suite.controller.registerHandler(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *UserControllerTestSuite) TestRegisterHandler_FailBinding() {
	recorder := httptest.NewRecorder()
	reqBody, _ := json.Marshal(UserProfileRequest.FullName)
	request, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(reqBody))
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	suite.controller.registerHandler(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *UserControllerTestSuite) TestUpdateHandler_Success() {
	suite.usecaseMock.On("Update", UserProfileRequest).Return(UserProfileRequest, nil)
	recorder := httptest.NewRecorder()
	reqBody, _ := json.Marshal(UserProfileRequest)
	request, _ := http.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(reqBody))
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	suite.controller.updateHandler(ctx)

	var response response.SingleResponse
	json.Unmarshal([]byte(recorder.Body.Bytes()), &response)

	assert.Equal(suite.T(), http.StatusCreated, recorder.Code)
	assert.NotEmpty(suite.T(), response.Data)
}

func (suite *UserControllerTestSuite) TestUpdateHandler_FailUsecase() {
	expectErr := errors.New("Failed")
	suite.usecaseMock.On("Update", UserProfileRequest).Return(request.UserProfileRequest{}, expectErr)
	recorder := httptest.NewRecorder()
	reqBody, _ := json.Marshal(UserProfileRequest)
	request, _ := http.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(reqBody))
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	suite.controller.updateHandler(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *UserControllerTestSuite) TestUpdateHandler_FailBinding() {
	recorder := httptest.NewRecorder()
	reqBody, _ := json.Marshal(UserProfileRequest.FullName)
	request, _ := http.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(reqBody))
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	suite.controller.updateHandler(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *UserControllerTestSuite) TestDeleteHandler_Success() {
    suite.usecaseMock.On("Delete", UserProfile1.Id).Return(nil)
    recorder := httptest.NewRecorder()
    request, _ := http.NewRequest(http.MethodDelete, "/user/"+UserProfile1.Id, nil)
    ctx, _ := gin.CreateTestContext(recorder)
    ctx.Request = request
	ctx.AddParam("id", UserProfile1.Id)

    suite.controller.deleteByIdHandler(ctx)

    assert.Equal(suite.T(), http.StatusCreated, recorder.Code)
}

func (suite *UserControllerTestSuite) TestDeleteHandler_Fail() {
	expectErr := errors.New("Failed")
    suite.usecaseMock.On("Delete", UserProfile1.Id).Return(expectErr)
    recorder := httptest.NewRecorder()
    request, _ := http.NewRequest(http.MethodDelete, "/user/"+UserProfile1.Id, nil)
    ctx, _ := gin.CreateTestContext(recorder)
    ctx.Request = request
	ctx.AddParam("id", UserProfile1.Id)
	
    suite.controller.deleteByIdHandler(ctx)

    assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.usecaseMock = new(UsecaseMock)
	suite.routerMock = gin.Default()
	suite.controller = NewUserController(suite.usecaseMock, suite.routerMock)
}

func TestCustomerUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
