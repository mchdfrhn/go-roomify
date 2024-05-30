package repository

import (
	"database/sql"
	"errors"
	"go-roomify/model"
	"go-roomify/model/dto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

type FacilityRepoTestSuite struct {
	suite.Suite
	repo    FacilityRepository
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
}

func (suite *FacilityRepoTestSuite) SetupTest() {
	mockDb, mockSql, _ := sqlmock.New()
	facilityRepo := NewFacilityRepository(mockDb)
	suite.mockDb = mockDb
	suite.mockSql = mockSql
	suite.repo = facilityRepo

}

func TestFacilityRepoTestSuite(t *testing.T) {
	suite.Run(t, new(FacilityRepoTestSuite))
}

func (suite *FacilityRepoTestSuite) TestFacilityCreate_Success() {
	suite.mockSql.ExpectExec("INSERT INTO mst_facility").WithArgs(
		"F001",
		"meja kerja",
		"R001",
	).WillReturnResult(sqlmock.NewResult(1, 1))
	err := suite.repo.InsertFacility(dummyFacility[0])
	suite.Assert().NoError(err)
}

func (suite *FacilityRepoTestSuite) TestFacilityCreate_Fail() {
	expectErr := errors.New("Failed")
	suite.mockSql.ExpectExec("INSERT INTO mst_facility").WithArgs(
		"F001",
		"meja kerja",
		"R001",
	).WillReturnError(expectErr)
	err := suite.repo.InsertFacility(dummyFacility[0])
	suite.Assert().Equal(expectErr, err)
}

func (suite *FacilityRepoTestSuite) TestUpdateFacility_Success() {
	suite.mockSql.ExpectExec("UPDATE mst_facility SET name = \\$1,room_id = \\$2 WHERE id = \\$3").WithArgs(
		"meja kerja",
		"R001",
		"F001",
	).WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.repo.UpdateFacility(dummyFacility[0])
	suite.Assert().NoError(err)
}

func (suite *FacilityRepoTestSuite) TestUpdateFacility_Fail() {
	expectErr := errors.New("Failed")
	suite.mockSql.ExpectExec("UPDATE mst_facility SET name = \\$1,room_id = \\$2 WHERE id = \\$3").WithArgs(
		"meja kerja",
		"R001",
		"F001",
	).WillReturnError(expectErr)

	err := suite.repo.UpdateFacility(dummyFacility[0])
	suite.Assert().Equal(expectErr, err)
}

func (suite *FacilityRepoTestSuite) TestDeleteFacility_Success() {
	suite.mockSql.ExpectExec("DELETE FROM mst_facility WHERE id = \\$1").WithArgs(
		"F001",
	).WillReturnResult(sqlmock.NewResult(1, 1))

	err := suite.repo.DeleteFacility(dummyFacility[0].Id)
	suite.Assert().NoError(err)
}

func (suite *FacilityRepoTestSuite) TestDeleteFacility_Fail() {
	expectErr := errors.New("Failed")
	suite.mockSql.ExpectExec("DELETE FROM mst_facility WHERE id = \\$1").WithArgs(
		"F001",
	).WillReturnError(expectErr)

	err := suite.repo.DeleteFacility(dummyFacility[0].Id)
	suite.Assert().Equal(expectErr, err)
}

func (suite *FacilityRepoTestSuite) TestGetFacilityById_Success() {
	row := sqlmock.NewRows([]string{"id", "name", "room_id"}).AddRow(
		"F001",
		"meja kerja",
		"R001",
	)

	suite.mockSql.ExpectQuery("SELECT id,name,room_id FROM mst_facility WHERE id = \\$1").WithArgs(
		"F001",
	).WillReturnRows(row)

	facility, err := suite.repo.GetFacilityById(dummyFacility[0].Id)
	suite.Assert().NoError(err)
	suite.Assert().Equal(dummyFacility[0], facility)
}

func (suite *FacilityRepoTestSuite) TestGetFacilityById_Fail() {
	expectErr := errors.New("Failed")
	suite.mockSql.ExpectQuery("SELECT id,name,room_id FROM mst_facility WHERE id = \\$1").WithArgs(
		"F001",
	).WillReturnError(expectErr)

	_, err := suite.repo.GetFacilityById(dummyFacility[0].Id)
	suite.Assert().Equal(expectErr, err)
}

func (suite *FacilityRepoTestSuite) TestGetPagingFacility_Success() {
	page := 1
	size := 2
	skip := (page - 1) * size
	totalRows := 3

	suite.mockSql.ExpectQuery("SELECT id,name,room_id FROM mst_facility LIMIT \\$1 OFFSET \\$2").WithArgs(size, skip).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "room_id"}).
		AddRow("F001", "meja kerja", "R001").
		AddRow("F002", "Kursi kerja", "R002"))

	suite.mockSql.ExpectQuery("SELECT COUNT\\(id\\) FROM mst_facility").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(totalRows))

	facilities, paging, err := suite.repo.GetPagingFacility(page, size)

	// Assertions
	suite.Assert().NoError(err)
	suite.Assert().Equal(dummyFacility[:2], facilities)
	suite.Assert().Equal(totalRows, paging.TotalRows)
	suite.Assert().Equal(page, paging.Page)
	suite.Assert().Equal(size, paging.RowsPerPage)
	suite.Assert().Equal(2, paging.TotalPages)
}

func (suite *FacilityRepoTestSuite) TestGetPagingFacility_Fail() {
	expectErr := errors.New("ERROR")
	page := 1
	size := 2
	skip := (page - 1) * size
	totalRows := 3

	suite.mockSql.ExpectQuery("SELECT id,name,room_id FROM mst_facility LIMIT \\$1 OFFSET \\$2").WithArgs(size, skip).WillReturnError(expectErr)

	suite.mockSql.ExpectQuery("SELECT COUNT\\(id\\) FROM mst_facility").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(totalRows))
	facilities, paging, err := suite.repo.GetPagingFacility(page, size)

	suite.Assert().Error(expectErr, err)
	suite.Assert().Nil(facilities)
	suite.Assert().Equal(dto.Paging{}, paging)
}
