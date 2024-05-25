package utils

import (
	"fmt"
	"go-roomify/model/dto"
	"strconv"
	"github.com/gin-gonic/gin"
)

// Limit, Offset, Order

func validateRequestQueryParams(c *gin.Context) (dto.RequestQueryParam,error) {

	// Validasi Page
	page, err := strconv.Atoi(c.DefaultQuery("page","1"))
	if err != nil || page <= 0{
		return dto.RequestQueryParam{},fmt.Errorf("Invalid Page NUmber")
	}

	// Validasi Limit
	limit, err := strconv.Atoi(c.DefaultQuery("limit","5"))
	if err != nil || limit <= 0{
		return dto.RequestQueryParam{},fmt.Errorf("Invalid Limit NUmber")
	}

	order := c.DefaultQuery("order","id")
	sort := c.DefaultQuery("sort","asc")

	return dto.RequestQueryParam {
		QueryParams: dto.QueryParams{
			Order : order,
			Sort : sort,
		},
		PaginationParam: dto.PaginationParam{
			Page : page,
			Limit: limit,
		},
	},nil
}