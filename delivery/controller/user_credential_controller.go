package controller

import (
	"go-roomify/middleware"
	"go-roomify/model"
	"go-roomify/model/dto/request"
	"go-roomify/model/dto/response"
	"go-roomify/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserCredentialController struct {
	uc             usecase.UserCredentialUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (self *UserCredentialController) loginHandler(ctx *gin.Context) {
	var login_payload request.UserCredentialRequest

	if err := ctx.ShouldBindJSON(&login_payload); err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	r_user_cr, err := self.uc.LoginUser(login_payload)

	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SendSingleResponseCreated(ctx, r_user_cr, "Login Success")
}

func (self *UserCredentialController) createNewHandler(ctx *gin.Context) {
	var new_user model.UserCredential

	if err := ctx.ShouldBindJSON(&new_user); err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	r_user_cr, err := self.uc.CreateNewUser(new_user)

	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SendSingleResponseCreated(
		ctx,
		r_user_cr,
		"Success Create New User")
}

func (self *UserCredentialController) getListHandler(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	rows_user, paging, err := self.uc.GetListUser(page, size)

	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var data []any
	data = append(data, rows_user)
	response.SendSinglePageResponse(ctx, data, "Success Get List User Credential", paging)
}

func (self *UserCredentialController) getByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	r_user, err := self.uc.GetUserById(id)

	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, "Invalid User Id")
		return
	}

	response.SendSingleResponseData(
		ctx,
		r_user,
		"Success Get User Credential Data")
}

func (self *UserCredentialController) updatePasswordHandler(ctx *gin.Context) {
	var new_user request.UserUpdatePasswordRequest

	if err := ctx.ShouldBindJSON(&new_user); err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := self.uc.UpdatePassword(new_user)

	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.SendSingleResponse(
		ctx,
		"Success Update User Password",
	)
}

func (self *UserCredentialController) deleteHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := self.uc.DeleteUser(id)

	if err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	response.SendSingleResponse(
		ctx,
		"Success Delete User Credential Data",
	)
}

func (self *UserCredentialController) Route() {
	router := self.rg.Group("/auth")

	router.POST("/login", self.loginHandler)

	router.GET("/user/:id", self.authMiddleware.RequireToken("admin"), self.getByIdHandler)
	router.POST("/user", self.authMiddleware.RequireToken("admin"), self.createNewHandler)
	router.PUT("/user", self.authMiddleware.RequireToken("admin"), self.updatePasswordHandler)
	router.DELETE("/user/:id", self.authMiddleware.RequireToken("admin"), self.deleteHandler)
	router.GET("/user", self.authMiddleware.RequireToken("admin"), self.getListHandler)
}

func NewUserCredentialController(uc usecase.UserCredentialUsecase, router *gin.Engine, auth_middleware middleware.AuthMiddleware) *UserCredentialController {
	return &UserCredentialController{
		uc:             uc,
		rg:             &router.RouterGroup,
		authMiddleware: auth_middleware,
	}
}
