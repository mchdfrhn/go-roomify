package response

import (
	"go-roomify/model"
	"github.com/gin-gonic/gin"
)

func SendDivisionsResponse(ctx *gin.Context, statusCode int, divisions []model.Division) {
	ctx.JSON(statusCode, gin.H{"divisions": divisions})
}

func SendDivisionResponse(ctx *gin.Context, statusCode int, division model.Division) {
	ctx.JSON(statusCode, gin.H{"division": division})
}

