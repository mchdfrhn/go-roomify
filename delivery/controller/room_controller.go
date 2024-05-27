package controller

import (
	"go-roomify/model/dto/request"
	"go-roomify/model/dto/response"
	"go-roomify/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type roomController struct{
	ru usecase.RoomUsecase
	rg *gin.RouterGroup
}

func( rc *roomController ) createRoomHandler(ctx *gin.Context){

	var roomRequest request.RoomRequest
	if err := ctx.ShouldBindJSON( &roomRequest ); err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	createdRoom, code, err := rc.ru.CreateRoom( roomRequest )
	if err != nil {
		response.SendSingleResponseError(
			ctx,
			code,
			err.Error(),
		)

		return
	}

	response.SendSingleResponseCreated(
		ctx,
		createdRoom,
		"Success Create Room",
	)
	
}

func( rc *roomController ) getRoomByIdOrNameHandler(ctx *gin.Context){

	roomIdOrName := ctx.Param("idOrName")

	createdRoom, code, err := rc.ru.GetRoomByIdOrName( roomIdOrName )
	if err != nil {
		response.SendSingleResponseError(
			ctx,
			code,
			err.Error(),
		)

		return
	}

	response.SendSingleResponseCreated(
		ctx,
		createdRoom,
		"Success Get data Room",
	)
	
}

func( rc *roomController ) Route(){
	group := rc.rg.Group("/room")
	group.POST("/", rc.createRoomHandler)
	group.GET("/:idOrName", rc.getRoomByIdOrNameHandler)
}

func NewRoomController( ru usecase.RoomUsecase, rg *gin.RouterGroup ) *roomController{
	return &roomController{
		ru: ru,
		rg: rg,
	}
}