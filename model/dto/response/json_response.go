package response

import (
	"go-roomify/model/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendSingleResponseCreated(ctx *gin.Context, data any, descriptionMsg string) {
	ctx.JSON(http.StatusOK, &SingleResponse{
		Status: Status{
			Code:        http.StatusCreated,
			Description: descriptionMsg,
		},
		Data: data,
	})
}

func SendSingleResponseData(ctx *gin.Context, data any, descriptionMsg string) {
	ctx.JSON(http.StatusOK, &SingleResponse{
		Status: Status{
			Code:        http.StatusOK,
			Description: descriptionMsg,
		},
		Data: data,
	})
}

func SendSinglePageResponse(ctx *gin.Context, data []any, descriptionMsg string, paging dto.Paging) {
	ctx.JSON(http.StatusOK, &PagedResponse{
		Status: Status{
			Code:        http.StatusOK,
			Description: descriptionMsg,
		},
		Data:   data,
		Paging: paging,
	})
}

func SendSingleResponse(ctx *gin.Context, descriptionMsg string) {
	ctx.JSON(http.StatusOK, &Status{
		Code:        http.StatusOK,
		Description: descriptionMsg,
	},
	)
}

func SendSingleResponseError(ctx *gin.Context, code int, errorMessage string) {
	ctx.AbortWithStatusJSON(code, &Status{
		Code:        code,
		Description: errorMessage,
	})
}
