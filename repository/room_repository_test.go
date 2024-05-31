package repository

import (
	"database/sql"
	"errors"
	"go-roomify/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

var rooms = []model.Room{
	{
		Id: "1",
		Name: "B01",
		RoomType: "meeting",
		Capacity: 20,
		IsAvailable: true,
	},
	{
		Id: "2",
		Name: "B02",
		RoomType: "meeting",
		Capacity: 20,
		IsAvailable: true,
	},
}

type roomRepoTestSuite struct{
	suite.Suite
	repo RoomRepository
	mockDb *sql.DB
	mockSql sqlmock.Sqlmock
}

func (suite *roomRepoTestSuite) SetupTest(){
	mockDb, mockSql, _ := sqlmock.New()
	roomRepo := NewRoomRepository( mockDb )

	suite.mockDb = mockDb
	suite.mockSql = mockSql
	suite.repo = roomRepo
}

func ( suite *roomRepoTestSuite ) TestRoomCreate_Success(){

	suite.mockSql.ExpectExec(
		"INSERT INTO mst_room",
	).WithArgs(
		"1",
		"B01",
		"meeting",
		20,
		true,
	).WillReturnResult( 
		sqlmock.NewResult( 1, 1 ),
	)

	err := suite.repo.CreateRoom( rooms[0] )

	suite.Assert().NoError( err )

}

func ( suite *roomRepoTestSuite ) TestRoomCreate_Fail(){

	ExpectErr := errors.New("ExecQuery ' INSERT INTO mst_room (id,name,roomtype,capacity,is_available) VALUES ($1,$2,$3,$4,$5)', arguments do not match: argument 3 expected [string - 20] does not match actual [int64 - 20]")

	suite.mockSql.ExpectExec(
		"INSERT INTO mst_room",
	).WithArgs(
		"1",
		"B01",
		"meeting",
		"20",
		true,
	).WillReturnError(
		ExpectErr,
	)

	err := suite.repo.CreateRoom( rooms[0] )

	suite.Assert().Equal( ExpectErr , err )
	
}

func TestRoomRepoTestSuite( t *testing.T ){
	suite.Run( t, new( roomRepoTestSuite ) )
}