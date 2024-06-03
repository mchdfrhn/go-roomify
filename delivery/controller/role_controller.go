package controller

import (
	"go-roomify/middleware"
	"go-roomify/model"
	"go-roomify/model/dto/response"
	"go-roomify/usecase"
	"go-roomify/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleController struct {
	uc             usecase.RoleUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
	// authMiddleware middleware.AuthMiddleware
}

func (rc *RoleController) registerHandler(ctx *gin.Context) {
	var newRole model.Role
	if err := ctx.ShouldBindJSON(&newRole); err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	newRole.Id = uuid.NewString()
	err := rc.uc.RegisterRole(newRole)
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
		newRole,
		"Success Register new Role",
	)
}

func (rc *RoleController) findAllPageHandler(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	roles, paging, err := rc.uc.FindAllPaging(page, size)
	if err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)
	}
	var data []any
	data = append(data, roles)
	response.SendSinglePageResponse(
		ctx,
		data,
		"Success Get List Role",
		paging,
	)
}

func (rc *RoleController) findByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	role, err := rc.uc.FindRoleById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Success Get Role By Id",
		"data":    role,
	})
}

func (rc *RoleController) deleteByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing id parameter",
		})
		return
	}
	err := rc.uc.DeleteRoleById(id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Success Delete Role By Id",
	})
}

func (rc *RoleController) updateHandler(ctx *gin.Context) {
	var role model.Role
	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	err := rc.uc.UpdateRole(role)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Success Update Role By Id",
		"data":    role,
	})
}

func (rc *RoleController) Route() {
	router := rc.rg.Group("/roles")
	// router.Use(rc.authMiddleware.RequireToken("ADMIN"))
	router.Use(rc.authMiddleware.RequireToken(utils.USER_ROLE_ADMIN))

	router.POST("", rc.registerHandler)
	router.GET("", rc.findAllPageHandler)
	router.GET("/:id", rc.findByIdHandler)
	router.DELETE("/:id", rc.deleteByIdHandler)
	router.PUT("", rc.updateHandler)
}

func NewRoleController(uc usecase.RoleUsecase, router *gin.Engine, auth_middleware middleware.AuthMiddleware) *RoleController {
	return &RoleController{
		uc:             uc,
		rg:             &router.RouterGroup,
		authMiddleware: auth_middleware,
	}
}
