package controller

import (
	"go-roomify/middleware"
	"go-roomify/model/dto/request"
	"go-roomify/model/dto/response"
	"go-roomify/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	uc             usecase.UserProfileUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (cc *UserController) findAllPageHandler(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	users, paging, err := cc.uc.GetList(page, size)
	if err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	var data []any
	data = append(data, users)
	response.SendSinglePageResponse(
		ctx,
		data,
		"Success Get List User",
		paging,
	)
}

func (cc *UserController) findByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := cc.uc.GetById(id)
	if err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	response.SendSingleResponseData(
		ctx,
		user,
		"Success Get User By Id",
	)
}

func (cc *UserController) findByUsernameHandler(ctx *gin.Context) {
	username := ctx.Param("username")
	user, err := cc.uc.GetByUsername(username)
	if err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	response.SendSingleResponseData(
		ctx,
		user,
		"Success Get User By Username",
	)
}

func (cc *UserController) registerHandler(ctx *gin.Context) {
	var newUser request.UserProfileRequest
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	newUser.Id = uuid.NewString()
	user, err := cc.uc.Create(newUser)
	if err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	response.SendSingleResponseCreated(
		ctx,
		user,
		"Success Register new User",
	)
}

func (cc *UserController) updateHandler(ctx *gin.Context) {
	var newUser request.UserProfileRequest
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	user, err := cc.uc.Update(newUser)
	if err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	response.SendSingleResponseData(
		ctx,
		user,
		"Success Update User By Id",
	)
}

func (cc *UserController) deleteByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	err := cc.uc.Delete(id)
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
		"Success Delete User By Id",
	)
}

func (cc *UserController) Route() {
	router := cc.rg.Group("/user")
	router.Use(cc.authMiddleware.RequireToken("admin"))

	router.GET("", cc.findAllPageHandler)
	router.GET("/:id", cc.findByIdHandler)
	router.GET("user/:username", cc.findByUsernameHandler)
	router.POST("", cc.registerHandler)
	router.PUT("", cc.updateHandler)
	router.DELETE("/:id", cc.deleteByIdHandler)
}

func NewUserController(uc usecase.UserProfileUsecase, router *gin.Engine, auth_middleware middleware.AuthMiddleware) *UserController {
	return &UserController{
		uc:             uc,
		rg:             &router.RouterGroup,
		authMiddleware: auth_middleware,
	}
}
