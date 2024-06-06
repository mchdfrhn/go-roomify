package controller

import (
	"go-roomify/middleware"
	"go-roomify/model"
	"go-roomify/model/dto/request"
	"go-roomify/model/dto/response"
	"go-roomify/usecase"
	"go-roomify/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DivisionController struct {
	uc             usecase.DivisionUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (c *DivisionController) GetDivisions(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	divisions, paging, err := c.uc.GetAllDivisions(page, size)

	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var division []any
	division = append(division, divisions)
	response.SendSinglePageResponse(
		ctx,
		division,
		"Success Get List Division",
		paging,
	)

}

func (c *DivisionController) GetDivisionById(ctx *gin.Context) {
	id := ctx.Param("id")
	division, err := c.uc.GetDivisionById(id)
	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusNotFound, "Division not found")
		return
	}
	response.SendDivisionResponse(ctx, http.StatusOK, division)
}

func (c *DivisionController) CreateDivision(ctx *gin.Context) {
	var req request.DivisionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// division := model.Division{
	// 	Id:   req.Id,
	// 	Name: req.Name,
	// }
	division, err := c.uc.CreateDivision(req)
	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SendSingleResponseCreated(ctx, division, "Division created successfully")
}

func (c *DivisionController) UpdateDivision(ctx *gin.Context) {
	var req request.DivisionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	division := model.Division{
		Id:   req.Id,
		Name: req.Name,
	}
	err := c.uc.UpdateDivision(division)
	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SendSingleResponseData(ctx, division, "Division updated successfully")
}

func (c *DivisionController) DeleteDivision(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.uc.DeleteDivision(id)
	if err != nil {
		response.SendSingleResponseError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.SendSingleResponse(ctx, "Division deleted successfully")
}

func (c *DivisionController) Route() {
	rg := c.rg.Group("/divisions")
	rg.Use(c.authMiddleware.RequireToken(utils.USER_ROLE_ADMIN))

	rg.GET("/", c.GetDivisions)
	rg.GET("/:id", c.GetDivisionById)
	rg.POST("/", c.CreateDivision)
	rg.PUT("/", c.UpdateDivision)
	rg.DELETE("/:id", c.DeleteDivision)
}

func NewDivisionController(uc usecase.DivisionUsecase, router *gin.Engine, auth_middleware middleware.AuthMiddleware) *DivisionController {
	return &DivisionController{
		uc:             uc,
		rg:             &router.RouterGroup,
		authMiddleware: auth_middleware,
	}
}
