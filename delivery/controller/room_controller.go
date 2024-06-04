package controller

import (
	"go-roomify/model"
	"go-roomify/model/dto/request"
	"go-roomify/model/dto/response"
	"go-roomify/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	ru usecase.RoomUsecase
	rg *gin.RouterGroup
}

func (rc *RoomController) CreateRoomHandler(ctx *gin.Context) {

	var roomRequest request.RoomRequest
	if err := ctx.ShouldBindJSON(&roomRequest); err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	createdRoom, code, err := rc.ru.CreateRoom(roomRequest)
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

func (rc *RoomController) GetAllRoomHandler(ctx *gin.Context) {

	paramPage := ctx.Query("page")
	paramSize := ctx.Query("size")

	room, paging, code, err := rc.ru.GetAllRoom(paramPage, paramSize)
	if err != nil {
		response.SendSingleResponseError(
			ctx,
			code,
			err.Error(),
		)

		return
	}

	response.SendSinglePageResponse(
		ctx,
		room,
		"Success Get All Room",
		paging,
	)

}

func (rc *RoomController) GetRoomByIdOrNameHandler(ctx *gin.Context) {

	roomIdOrName := ctx.Param("idOrName")

	createdRoom, code, err := rc.ru.GetRoomByIdOrName(roomIdOrName)
	if err != nil {
		response.SendSingleResponseError(
			ctx,
			code,
			err.Error(),
		)

		return
	}

	response.SendSingleResponse(
		ctx,
		createdRoom,
		"Success Get data Room",
	)

}

func (rc *RoomController) UpdateRoomByIdHandler(ctx *gin.Context) {

	var updateRoom model.Room
	if err := ctx.ShouldBindJSON(&updateRoom); err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	updatedRoom, code, err := rc.ru.UpdateRoomById(updateRoom)
	if err != nil {
		response.SendSingleResponseError(
			ctx,
			code,
			err.Error(),
		)

		return
	}

	response.SendSingleResponse(
		ctx,
		updatedRoom,
		"Success Update Room",
	)

}

func (rc *RoomController) DeleteRoomByIdHandler(ctx *gin.Context) {

	roomId := ctx.Param("id")

	code, err := rc.ru.DeleteRooomById(roomId)
	if err != nil {
		response.SendSingleResponseError(
			ctx,
			code,
			err.Error(),
		)

		return
	}

	response.SendSingleResponse(
		ctx,
		code,
		"Success delete data Room",
	)

}

func (rc *RoomController) GetRoomAvailableHandler(ctx *gin.Context) {

	availableRoom, code, err := rc.ru.GetAvailableRoom()
	if err != nil {
		response.SendSingleResponseError(
			ctx,
			code,
			err.Error(),
		)

		return
	}

	response.SendSingleResponse(
		ctx,
		availableRoom,
		"Success get available rooms",
	)

}

func (rc *RoomController) Route() {
	group := rc.rg.Group("/room")
	group.POST("/", rc.CreateRoomHandler)
	group.GET("/", rc.GetAllRoomHandler)
	group.GET("/:idOrName", rc.GetRoomByIdOrNameHandler)
	group.PUT("/", rc.UpdateRoomByIdHandler)
	group.DELETE("/:id", rc.DeleteRoomByIdHandler)
	group.GET("/available", rc.GetRoomAvailableHandler)
}

func NewRoomController(ru usecase.RoomUsecase, rg *gin.Engine) *RoomController {
	return &RoomController{
		ru: ru,
		rg: &rg.RouterGroup,
	}
}
