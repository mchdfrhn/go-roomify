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
	paramType := ctx.Query("type")

	room, paging, code, err := rc.ru.GetAllRoom(paramPage, paramSize, paramType)
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

	response.SendSingleResponseData(
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

	response.SendSingleResponseData(
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

	response.SendSingleResponseData(
		ctx,
		updatedRoom,
		"Success Update Room",
	)
}

func (rc *roomController) GetRoomAvailableHandler(ctx *gin.Context) {
	paramType := ctx.Query("type")

	availableRoom, code, err := rc.ru.GetAvailableRoom(paramType)
	if err != nil {
		response.SendSingleResponseError(
			ctx,
			code,
			err.Error(),
		)

		return
	}

	response.SendSingleResponseData(
		ctx,
		availableRoom,
		"Success get available rooms",
	)
}

func (rc *roomController) Route() {
	group := rc.rg.Group("/room")

	admin_Middleware := rc.authMiddleware.RequireToken(utils.USER_ROLE_ADMIN)
	group.POST("/", admin_Middleware, rc.createRoomHandler)
	group.PUT("/", admin_Middleware, rc.updateRoomByIdHandler)
	group.DELETE("/:id", admin_Middleware, rc.deleteRoomByIdHandler)

	admin_GA_Middleware := rc.authMiddleware.RequireToken(utils.USER_ROLE_ADMIN, utils.USER_ROLE_GA)
	group.PUT("/status", admin_GA_Middleware, rc.updateRoomByIdIsAvailableHandler)

	admin_GA_Employee_Middleware := rc.authMiddleware.RequireToken(utils.USER_ROLE_ADMIN, utils.USER_ROLE_GA, utils.USER_ROLE_EMPLOYEE)
	group.GET("/", admin_GA_Employee_Middleware, rc.getAllroomHandler)
	group.GET("/:idOrName", admin_GA_Employee_Middleware, rc.getRoomByIdOrNameHandler)
	group.GET("/available", admin_GA_Employee_Middleware, rc.GetRoomAvailableHandler)
}

func NewRoomController(ru usecase.RoomUsecase, router_group *gin.RouterGroup, auth_middleware middleware.AuthMiddleware) *roomController {
	return &roomController{
		ru:             ru,
		rg:             router_group,
		authMiddleware: auth_middleware,
	}
}
