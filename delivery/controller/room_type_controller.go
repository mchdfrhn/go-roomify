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

type roomTypeController struct {
	rtu            usecase.RoomTypeUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (rtc *roomTypeController) createRoomTypeHandler(ctx *gin.Context) {
	var roomTypeRequest request.RoomTypeRequest
	if err := ctx.ShouldBindJSON(&roomTypeRequest); err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	createdRoomType, code, err := rtc.rtu.CreateRoomType(roomTypeRequest)
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
		createdRoomType,
		"Success create roomtype",
	)
}

func (rtc *roomTypeController) getAllRoomTypeHandler(ctx *gin.Context) {

	roomType, code, err := rtc.rtu.GetAllRoomType()
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
		roomType,
		"Success get all Roomtype",
	)
}

func (rtc *roomTypeController) getRoomTypeByIdOrNameHandler(ctx *gin.Context) {
	roomTypeIdOrName := ctx.Param("idOrName")

	getRoomType, code, err := rtc.rtu.GetRoomTypeByIdOrName(roomTypeIdOrName)
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
		getRoomType,
		"Success Get data roomtype by id or name",
	)
}

func (rtc *roomTypeController) updateRoomTypeByIdHandler(ctx *gin.Context) {
	var updateRoomType model.RoomType

	if err := ctx.ShouldBindJSON(&updateRoomType); err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	updatedRoomType, code, err := rtc.rtu.UpdateRoomTypeById(updateRoomType)

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
		updatedRoomType,
		"Success update roomtype",
	)
}

func (rtc *roomTypeController) deleteRoomTypeByIdHandler(ctx *gin.Context) {
	roomTypeId := ctx.Param("id")

	code, err := rtc.rtu.DeleteRoomTypeById(roomTypeId)
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
		"Success delete data roomtype",
	)
}

func (rtc *roomTypeController) Route() {
	group := rtc.rg.Group("/type/room")

	group.Use(rtc.authMiddleware.RequireToken(utils.USER_ROLE_ADMIN, utils.USER_ROLE_GA))

	group.GET("/", rtc.getAllRoomTypeHandler)
	group.GET("/:idOrName", rtc.getRoomTypeByIdOrNameHandler)
	group.POST("/", rtc.createRoomTypeHandler)
	group.PUT("/", rtc.updateRoomTypeByIdHandler)
	group.DELETE("/:id", rtc.deleteRoomTypeByIdHandler)
}

func NewRoomTypeController(rtu usecase.RoomTypeUsecase, router_group *gin.RouterGroup, auth_middleware middleware.AuthMiddleware) *roomTypeController {
	return &roomTypeController{
		rtu:            rtu,
		rg:             router_group,
		authMiddleware: auth_middleware,
	}
}
