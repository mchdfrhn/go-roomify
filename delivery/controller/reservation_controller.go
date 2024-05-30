package controller

import (
	"fmt"
	"go-roomify/middleware"
	"go-roomify/model"
	"go-roomify/model/dto/request"
	"go-roomify/model/dto/response"
	"go-roomify/usecase"
	"go-roomify/utils"
	"net/http"

	//"strconv"

	"github.com/gin-gonic/gin"
)

type ReservationController struct {
	uc             usecase.ReservationUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (self *ReservationController) createNewHandler(ctx *gin.Context) {
	var new_reservation model.Reservation

	if err := ctx.ShouldBindJSON(&new_reservation); err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	r_reservation, err := self.uc.CreateRequest(new_reservation)

	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SendSingleResponseCreated(ctx, r_reservation, fmt.Sprintf("Success Create New Reservation Request"))
}

func (self *ReservationController) changeStatus(ctx *gin.Context, reserv_status int, success_message string) {
	var new_reservation_status request.ReservationStatus

	if err := ctx.ShouldBindJSON(&new_reservation_status); err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	new_reservation_status.Status = reserv_status
	err := self.uc.ChangeStatus(new_reservation_status)

	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SendSingleResponse(ctx, nil, success_message)
}

func (self *ReservationController) cancelHandler(ctx *gin.Context) {
	self.changeStatus(ctx, utils.RESERV_STATUS_CANCEL, "Success Cancel Room Reservation Request")
}

func (self *ReservationController) acceptHandler(ctx *gin.Context) {
	self.changeStatus(ctx, utils.RESERV_STATUS_ACCEPTED, "Success Accept Room Reservation Request")
}

func (self *ReservationController) declineHandler(ctx *gin.Context) {
	self.changeStatus(ctx, utils.RESERV_STATUS_DECLINE, "Success Decline Room Reservation Request")
}

func (self *ReservationController) Route() {
	router := self.rg.Group("/reservation")
	// router.Use(self.authMiddleware.RequireToken("admin", "employee", "ga"))
	router.POST("", self.createNewHandler, self.authMiddleware.RequireToken("admin", "employee", "ga"))
	router.PUT("/cancel", self.cancelHandler, self.authMiddleware.RequireToken("employee"))
	router.PUT("/accept", self.acceptHandler, self.authMiddleware.RequireToken("admin", "ga"))
	router.PUT("/decline", self.declineHandler, self.authMiddleware.RequireToken("admin", "ga"))
}

func NewReservationController(uc usecase.ReservationUsecase, router *gin.Engine, auth_middleware middleware.AuthMiddleware) *ReservationController {
	return &ReservationController{
		uc:             uc,
		rg:             &router.RouterGroup,
		authMiddleware: auth_middleware,
	}
}
