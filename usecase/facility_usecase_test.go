package usecase

import (
	"errors"
	"go-roomify/model"
	"go-roomify/model/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyFacility = []model.Facility{
	{
		Id:     "F001",
		Name:   "meja kerja",
		RoomId: "R001",
	},
	{
		Id:     "F002",
		Name:   "Kursi kerja",
		RoomId: "R002",
	},
	{
		Id:     "F003",
		Name:   "Proyektor",
		RoomId: "R003",
	},
}

type RepoMock struct {
	mock.Mock
}

type FacilityUsecaseTestSuite struct {
	suite.Suite
	repoMock *RepoMock
	usecase  FacilityUsecase
}

func (suite *FacilityUsecaseTestSuite) SetupTest() {
	suite.repoMock = new(RepoMock)
	suite.usecase = NewFacilityUseCase(suite.repoMock)
}

func TestFacilityUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(FacilityUsecaseTestSuite))
}

func (r *RepoMock) GetPagingFacility(page int, size int) ([]model.Facility, dto.Paging, error) {
	args := r.Called(page, size)
	if args.Get(2) != nil {
		return []model.Facility{}, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.Facility), dto.Paging{}, nil
}

// REPOMOCK
func (r *RepoMock) GetFacilityById(id string) (model.Facility, error) {
	args := r.Called(id)
	if args.Get(1) != nil {
		return model.Facility{}, args.Error(1)
	}
	return args.Get(0).(model.Facility), nil
}

func (r *RepoMock) InsertFacility(newFacility model.Facility) error {
	args := r.Called(newFacility)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (r *RepoMock) UpdateFacility(newFacility model.Facility) error {
	args := r.Called(newFacility)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (r *RepoMock) DeleteFacility(id string) error {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

// TESTING FACILITY USECASE

func (suite *FacilityUsecaseTestSuite) TestFacilityPagingList_Success() {
	suite.repoMock.On("GetPagingFacility", 10, 0).Return(dummyFacility, dto.Paging{}, nil)
	actual, _, err := suite.usecase.FindAllPagingFacility(10, 0)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dummyFacility, actual)
}

func (suite *FacilityUsecaseTestSuite) TestFacilityPagingList_Fail() {
	expectErr := errors.New("Failed")
	suite.repoMock.On("GetPagingFacility", 10, 0).Return([]model.Facility{}, dto.Paging{}, expectErr)
	_, _, err := suite.usecase.FindAllPagingFacility(10, 0)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectErr, err)
}

func (suite *FacilityUsecaseTestSuite) TestFindFacilityById_Success() {
	suite.repoMock.On("GetFacilityById", dummyFacility[0].Id).Return(dummyFacility[0], nil)
	actual, err := suite.usecase.FindFacilityById(dummyFacility[0].Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), dummyFacility[0], actual)
}

func (suite *FacilityUsecaseTestSuite) TestUserFindFacilityById_Fail() {
	expectErr := errors.New("Failed")
	suite.repoMock.On("GetFacilityById", dummyFacility[0].Id).Return(model.Facility{}, expectErr)
	_, err := suite.usecase.FindFacilityById(dummyFacility[0].Id)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectErr, err)
}

func (suite *FacilityUsecaseTestSuite) TestInputFacility_Success() {
	suite.repoMock.On("InsertFacility", dummyFacility[0]).Return(nil)
	err := suite.usecase.InputFacility(dummyFacility[0])
	assert.NoError(suite.T(), err)
}

func (suite *FacilityUsecaseTestSuite) TestInputFacility_Fail() {
	suite.repoMock.On("InsertFacility", dummyFacility[0]).Return(errors.New("Failed"))
	err := suite.usecase.InputFacility(dummyFacility[0])
	assert.Error(suite.T(), err)
}

func (suite *FacilityUsecaseTestSuite) TestUpdatedFacility_Success() {
	suite.repoMock.On("UpdateFacility", dummyFacility[0]).Return(nil)
	err := suite.usecase.UpdatedFacility(dummyFacility[0])
	assert.NoError(suite.T(), err)
}

func (suite *FacilityUsecaseTestSuite) TestUpdatedFacility_Fail() {
	expectErr := errors.New("Failed")
	suite.repoMock.On("UpdateFacility", dummyFacility[0]).Return(expectErr)
	err := suite.usecase.UpdatedFacility(dummyFacility[0])
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectErr, err)
}

func (suite *FacilityUsecaseTestSuite) TestDeletedFacility_Success() {
	suite.repoMock.On("DeleteFacility", dummyFacility[0].Id).Return(nil)
	err := suite.usecase.DeletedFacility(dummyFacility[0].Id)
	assert.Nil(suite.T(), err)
}

func (suite *FacilityUsecaseTestSuite) TestDeletedFacility_Error() {
	expectErr := errors.New("Failed")
	suite.repoMock.On("DeleteFacility", dummyFacility[0].Id).Return(expectErr)
	err := suite.usecase.DeletedFacility(dummyFacility[0].Id)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectErr, err)
}
