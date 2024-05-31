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

	"github.com/dgrijalva/jwt-go"
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

	jwt_claims := ctx.MustGet("claims").(jwt.MapClaims)

	new_reservation.UserProfileId = jwt_claims["user_id"].(string)

	r_reservation, err := self.uc.CreateRequest(new_reservation)

	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SendSingleResponseCreated(ctx, r_reservation, fmt.Sprintf("Success Create New Reservation Request"))
}

func (self *ReservationController) statusHandler(ctx *gin.Context) {
	var new_reserv_status request.ReservationStatusRequest

	if err := ctx.ShouldBindJSON(&new_reserv_status); err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	jwt_claims := ctx.MustGet("claims").(jwt.MapClaims)

	//user_id := jwt_claims["user_id"]
	user_role := jwt_claims["role"]

	if user_role == "employee" && new_reserv_status.StatusId != utils.RESERV_STATUS_CANCEL {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, "Invalid Status Code")
		return
	}

	if user_role != "employee" && new_reserv_status.StatusId == utils.RESERV_STATUS_CANCEL {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, "Invalid Status Code")
		return
	}

	err := self.uc.ChangeStatus(new_reserv_status)

	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SendSingleResponse(ctx, nil, "Success Change Reservation Status")
}

func (self *ReservationController) getListByTokenHandler(ctx *gin.Context) {
	jwt_claims := ctx.MustGet("claims").(jwt.MapClaims)

	fl_reserv_get_list := request.ReservationGetListFilter{
		UserId:       jwt_claims["user_id"].(string),
		UserRole:     jwt_claims["role"].(string),
		FilterStatus: ctx.DefaultQuery("fl_status", ""),
	}

	fmt.Println(fl_reserv_get_list)

	r_reservation, err := self.uc.GetListByToken(fl_reserv_get_list)

	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var data []any
	data = append(data, r_reservation)

	response.SendSingleResponse(ctx, data, "Success Get List Reservation")
}

func (self *ReservationController) Route() {
	router := self.rg.Group("/reservation")
	router.Use(self.authMiddleware.RequireToken("admin", "employee", "ga"))
	// router.POST("/register", self.createNewHandler, self.authMiddleware.RequireToken("admin"))
	// router.POST("/delete", self.createNewHandler, self.authMiddleware.RequireToken("admin"))

	router.POST("", self.createNewHandler)
	router.GET("", self.getListByTokenHandler)
	router.PUT("/status", self.statusHandler)
	// router.PUT("/accept", self.acceptHandler, self.authMiddleware.RequireToken("admin", "ga"))
	// router.PUT("/decline", self.declineHandler, self.authMiddleware.RequireToken("admin", "ga"))
}

func NewReservationController(uc usecase.ReservationUsecase, router *gin.Engine, auth_middleware middleware.AuthMiddleware) *ReservationController {
	return &ReservationController{
		uc:             uc,
		rg:             &router.RouterGroup,
		authMiddleware: auth_middleware,
	}
}
