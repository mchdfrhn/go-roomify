package controller

import (
	"go-roomify/model"
	//"go-roomify/model/dto/request"
	"go-roomify/middleware"
	"go-roomify/model/dto/response"
	"go-roomify/usecase"
	"net/http"
	"fmt"
	//"strconv"

	"github.com/gin-gonic/gin"
)

type ReservationController struct {
	uc usecase.ReservationUsecase
	rg *gin.RouterGroup
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

func (self *ReservationController) cancelHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := self.uc.CancelRequest(id)

	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, response.Status{
		Code: http.StatusOK,
		Description: "Success Cancel Room Reservation Request",
	})
}

func (self *ReservationController) Route() {
	router := self.rg.Group("/reservation")
	router.Use(self.authMiddleware.RequireToken("admin", "employee", "ga"))
	router.POST("", self.createNewHandler)
	router.PUT("/cancel/:id", self.cancelHandler)
}

func NewReservationController(uc usecase.ReservationUsecase, router *gin.Engine, auth_middleware middleware.AuthMiddleware,) *ReservationController {
	return &ReservationController{
		uc: uc,
		rg : &router.RouterGroup,
		authMiddleware: auth_middleware,
	}
}
