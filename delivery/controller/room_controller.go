package controller

import (
	"go-roomify/middleware"
	"go-roomify/model"
	"go-roomify/model/dto/request"
	"go-roomify/model/dto/response"
	"go-roomify/usecase"
	"go-roomify/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type roomController struct {
	ru             usecase.RoomUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (rc *roomController) createRoomHandler(ctx *gin.Context) {

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

func (rc *roomController) getAllroomHandler(ctx *gin.Context) {

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

func (rc *roomController) getRoomByIdOrNameHandler(ctx *gin.Context) {

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

func (rc *roomController) updateRoomByIdHandler(ctx *gin.Context) {
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

func (rc *roomController) deleteRoomByIdHandler(ctx *gin.Context) {

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

func (rc *roomController) updateRoomByIdIsAvailableHandler(ctx *gin.Context) {
	var updateRoom request.RoomStatusRequest

	if err := ctx.ShouldBindJSON(&updateRoom); err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	updatedRoom, code, err := rc.ru.UpdateRoomByIdIsAvailableOnly(updateRoom.Id, *updateRoom.IsAvailable)

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

func (rc *roomController) GetRoomAvailableHandler(ctx *gin.Context) {

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

func (rc *roomController) Route() {
	group := rc.rg.Group("/room")
	//group.Use()

	group.POST("/", rc.authMiddleware.RequireToken(utils.USER_ROLE_ADMIN), rc.createRoomHandler)
	group.GET("/",
		rc.authMiddleware.RequireToken(
			utils.USER_ROLE_ADMIN,
			utils.USER_ROLE_GA,
			utils.USER_ROLE_EMPLOYEE),
		rc.getAllroomHandler)
	group.GET("/:idOrName",
		rc.authMiddleware.RequireToken(
			utils.USER_ROLE_ADMIN,
			utils.USER_ROLE_GA,
			utils.USER_ROLE_EMPLOYEE),
		rc.getRoomByIdOrNameHandler)
	group.PUT("/", rc.authMiddleware.RequireToken(utils.USER_ROLE_ADMIN), rc.updateRoomByIdHandler)
	group.PUT("/status", rc.authMiddleware.RequireToken(
		utils.USER_ROLE_ADMIN, utils.USER_ROLE_GA), rc.updateRoomByIdIsAvailableHandler)
	group.DELETE("/:id", rc.authMiddleware.RequireToken(
		utils.USER_ROLE_ADMIN), rc.deleteRoomByIdHandler)
	group.GET("/available",
		rc.authMiddleware.RequireToken(
			utils.USER_ROLE_ADMIN,
			utils.USER_ROLE_GA,
			utils.USER_ROLE_EMPLOYEE),
		rc.GetRoomAvailableHandler)
}

func NewRoomController(ru usecase.RoomUsecase, rg *gin.Engine, auth_middleware middleware.AuthMiddleware) *roomController {
	return &roomController{
		ru:             ru,
		rg:             &rg.RouterGroup,
		authMiddleware: auth_middleware,
	}
}
