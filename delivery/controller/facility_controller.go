package controller

import (
	"go-roomify/model"
	"go-roomify/model/dto/response"
	"go-roomify/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FacilityController struct {
	uf usecase.FacilityUsecase
	rg *gin.RouterGroup
}

func NewFacilityController(uf usecase.FacilityUsecase, router *gin.Engine) *FacilityController {
	return &FacilityController{
		uf: uf,
		rg: &router.RouterGroup,
	}
}

func (fc *FacilityController) FindAllPagingHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	facilities, paging, err := fc.uf.FindAllPagingFacility(page, size)
	if err != nil {
		response.SendSingleResponseError(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	var data []any
	data = append(data, facilities)
	response.SendSinglePageResponse(
		c,
		data,
		"Success Get List Facility",
		paging,
	)
}

func (fc *FacilityController) FindByIdHandler(c *gin.Context) {
	id := c.Param("id")
	facilities, err := fc.uf.FindFacilityById(id)
	if err != nil {
		response.SendSingleResponseError(
			c,
			http.StatusBadRequest,
			"Invalid Id",
		)
		return
	}
	response.SendSingleResponseCreated(
		c,
		facilities,
		"Success Get Id Facility",
	)
}

func (fc *FacilityController) InsertHandler(c *gin.Context) {
	var newFacility model.Facility
	if err := c.ShouldBindJSON(&newFacility); err != nil {
		response.SendSingleResponseError(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	newFacility.Id = uuid.NewString()
	err := fc.uf.InputFacility(newFacility)
	if err != nil {
		response.SendSingleResponseError(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	response.SendSingleResponseCreated(
		c,
		newFacility,
		"Success Register new Facility",
	)
}

func (fc *FacilityController) UpdateHandler(c *gin.Context) {
	var facilities model.Facility
	if err := c.ShouldBindJSON(&facilities); err != nil {
		response.SendSingleResponseError(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	err := fc.uf.UpdatedFacility(facilities)
	if err != nil {
		response.SendSingleResponseError(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	response.SendSingleResponse(
		c,
		facilities,
		"Success Facility Update",
	)
}

func (fc *FacilityController) DeleteHandler(c *gin.Context) {
	id := c.Param("id")
	err := fc.uf.DeletedFacility(id)
	if err != nil {
		response.SendSingleResponseError(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}
	response.SendSingleResponseCreated(
		c,
		http.StatusOK,
		"Success Facility Delete",
	)
}

func (fc *FacilityController) Route() {
	router := fc.rg.Group("/facility")
	router.GET("", fc.FindAllPagingHandler)
	router.GET("/:id", fc.FindByIdHandler)
	router.POST("", fc.InsertHandler)
	router.PUT("", fc.UpdateHandler)
	router.DELETE("/:id", fc.DeleteHandler)
}
