package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-roomify/model"
	"go-roomify/model/dto"
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

var dummyFacility = []model.Facility{
	{
		Id:     "538c7f96-b164-4f1b-97bb-9f4bb472e89f",
		Name:   "meja kerja",
		RoomId: "750e8400-e29b-41d4-a716-446655440021",
	},
	{
		Id:     "750e8400-e29b-41d4-a716-446655440002",
		Name:   "Kursi kerja",
		RoomId: "750e8400-e29b-41d4-a716-446655440022",
	},
	{
		Id:     "750e8400-e29b-41d4-a716-446655440003",
		Name:   "Proyektor",
		RoomId: "750e8400-e29b-41d4-a716-446655440023",
	},
}

type UsecaseMock struct {
	mock.Mock
}

type FacilityControllerTestSuite struct {
	suite.Suite
	routerMock  *gin.Engine
	usecaseMock *UsecaseMock
	controller  *FacilityController
}

func (u *UsecaseMock) FindAllPagingFacility(page int, size int) ([]model.Facility, dto.Paging, error) {
	args := u.Called()
	if args.Get(2) != nil {
		return []model.Facility{}, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.Facility), dto.Paging{}, nil
}

func (u *UsecaseMock) FindFacilityById(id string) (model.Facility, error) {
	args := u.Called(id)
	if args.Get(1) != nil {
		return model.Facility{}, args.Error(1)
	}
	return args.Get(0).(model.Facility), nil
}

func (u *UsecaseMock) InputFacility(newFacility model.Facility) error {
	args := u.Called(newFacility)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (u *UsecaseMock) UpdatedFacility(newFacility model.Facility) error {
	args := u.Called(newFacility)
	if args.Get(1) != nil {
		return args.Error(1)
	}
	return nil
}

func (u *UsecaseMock) DeletedFacility(id string) error {
	args := u.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (suite *FacilityControllerTestSuite) TestGetPagingHandler_Success() {
	suite.usecaseMock.On("FindAllPagingFacility").Return(dummyFacility, dto.Paging{}, nil)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/facility", nil)
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	suite.controller.FindAllPagingHandler(ctx)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *FacilityControllerTestSuite) TestFindAllPagingFacilitytHandler_Fail() {
	expectErr := errors.New("Failed")
	suite.usecaseMock.On("FindAllPagingFacility").Return([]model.Facility{}, dto.Paging{}, expectErr)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/facility", nil)
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	suite.controller.FindAllPagingHandler(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *FacilityControllerTestSuite) TestFindFacilityByIdHandler_Success() {
	suite.usecaseMock.On("FindFacilityById", dummyFacility[0].Id).Return(dummyFacility[0], nil)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/facility/"+dummyFacility[0].Id, nil)
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request
	ctx.AddParam("id", dummyFacility[0].Id)

	suite.controller.FindByIdHandler(ctx)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *FacilityControllerTestSuite) TestGetFacilityByIdHandler_Fail() {
	expectErr := errors.New("Failed")
	suite.usecaseMock.On("FindFacilityById", dummyFacility[0].Id).Return(model.Facility{}, expectErr)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/facility/"+dummyFacility[0].Id, nil)
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request
	ctx.AddParam("id", dummyFacility[0].Id)

	suite.controller.FindByIdHandler(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *FacilityControllerTestSuite) TestInputFacilityHandler_Success() {
	my_rand := rand.New(rand.NewSource(42))
	uuid.SetRand(my_rand)
	dummyFacility[0].Id = uuid.NewString()

	my_rand2 := rand.New(rand.NewSource(42))
	uuid.SetRand(my_rand2)

	suite.usecaseMock.On("InputFacility", dummyFacility[0]).Return(nil)
	recorder := httptest.NewRecorder()
	reqBody, _ := json.Marshal(dummyFacility[0])
	request, _ := http.NewRequest(http.MethodPost, "/facility", bytes.NewBuffer(reqBody))
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	suite.controller.InsertHandler(ctx)

	var response response.SingleResponse
	json.Unmarshal([]byte(recorder.Body.Bytes()), &response)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.NotEmpty(suite.T(), response.Data)
}

func (suite *FacilityControllerTestSuite) TestInputFacilityHandler_FailUsecase() {
	my_rand := rand.New(rand.NewSource(42))
	uuid.SetRand(my_rand)
	dummyFacility[0].Id = uuid.NewString()

	my_rand2 := rand.New(rand.NewSource(42))
	uuid.SetRand(my_rand2)

	expectErr := errors.New("Failed")
	suite.usecaseMock.On("InputFacility", dummyFacility[0]).Return(expectErr)
	recorder := httptest.NewRecorder()
	reqBody, _ := json.Marshal(dummyFacility[0])
	request, _ := http.NewRequest(http.MethodPost, "/facility", bytes.NewBuffer(reqBody))
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	suite.controller.InsertHandler(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *FacilityControllerTestSuite) TestInputHandler_FailBinding() {
	recorder := httptest.NewRecorder()
	reqBody, _ := json.Marshal(dummyFacility[0].Name)
	request, _ := http.NewRequest(http.MethodPost, "/facility", bytes.NewBuffer(reqBody))
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	suite.controller.InsertHandler(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *FacilityControllerTestSuite) TestUpdatedHandler_Success() {
	suite.usecaseMock.On("UpdatedFacility", dummyFacility[0]).Return(dummyFacility[0], nil)
	recorder := httptest.NewRecorder()
	reqBody, _ := json.Marshal(dummyFacility[0])
	request, _ := http.NewRequest(http.MethodPut, "/facility", bytes.NewBuffer(reqBody))
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	suite.controller.UpdateHandler(ctx)

	var response response.SingleResponse
	json.Unmarshal([]byte(recorder.Body.Bytes()), &response)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.NotEmpty(suite.T(), response.Data)
}

func (suite *FacilityControllerTestSuite) TestUpdatedHandler_FailUsecase() {
	suite.usecaseMock.On("UpdatedFacility", dummyFacility[0]).Return(nil, errors.New("Failed to update facility"))
	recorder := httptest.NewRecorder()
	reqBody, _ := json.Marshal(dummyFacility[0].Name)
	request, _ := http.NewRequest(http.MethodPut, "/facility", bytes.NewBuffer(reqBody))
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	suite.controller.UpdateHandler(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)

}

func (suite *FacilityControllerTestSuite) TestUpdatedHandler_FailBinding() {
	recorder := httptest.NewRecorder()
	reqBody, _ := json.Marshal(dummyFacility[0].Name)
	request, _ := http.NewRequest(http.MethodPut, "/facility", bytes.NewBuffer(reqBody))
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	suite.controller.UpdateHandler(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *FacilityControllerTestSuite) TestDeletedHandler_Success() {
	suite.usecaseMock.On("DeletedFacility", dummyFacility[0].Id).Return(nil)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, "/facility/"+dummyFacility[0].Id, nil)
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request
	ctx.AddParam("id", dummyFacility[0].Id)

	suite.controller.DeleteHandler(ctx)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
}

func (suite *FacilityControllerTestSuite) TestDeletedHandler_Fail() {
	expectErr := errors.New("Failed")
	suite.usecaseMock.On("DeletedFacility", dummyFacility[0].Id).Return(expectErr)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodDelete, "/user/"+dummyFacility[0].Id, nil)
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request
	ctx.AddParam("id", dummyFacility[0].Id)

	suite.controller.DeleteHandler(ctx)

	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

func (suite *FacilityControllerTestSuite) SetupTest() {
	suite.usecaseMock = new(UsecaseMock)
	suite.routerMock = gin.Default()
	suite.controller = NewFacilityController(suite.usecaseMock, suite.routerMock)
}

func TestCustomerUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(FacilityControllerTestSuite))
}
