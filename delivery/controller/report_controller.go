package controller

import (
	"go-roomify/middleware"
	"go-roomify/model/dto/response"
	"go-roomify/usecase"
	"go-roomify/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/sqltocsv"
)

type reportController struct {
	ru             usecase.ReportUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func (rc *reportController) downloadReportByYearHandler(ctx *gin.Context) {

	ctx.Header("Content-type", "text/csv")
	ctx.Header("Content-Disposition", "attachment; filename=\"report.csv\"")

	startYear := ctx.Query("s")
	endYear := ctx.Query("e")

	rows, code, err := rc.ru.DownloadReportByYear(
		startYear,
		endYear,
	)

	if err != nil {
		response.SendSingleResponseError(
			ctx,
			code,
			err.Error(),
		)

		return
	}

	err = sqltocsv.Write(ctx.Writer, rows)
	if err != nil {
		response.SendSingleResponseError(
			ctx,
			http.StatusInternalServerError,
			err.Error(),
		)

		return
	}

}

func (rc *reportController) Route() {
	group := rc.rg.Group("/report")
	group.Use(rc.authMiddleware.RequireToken(utils.USER_ROLE_ADMIN, utils.USER_ROLE_GA))

	group.GET("/", rc.downloadReportByYearHandler)
}

func NewReportController(ru usecase.ReportUsecase, rg *gin.Engine, auth_middleware middleware.AuthMiddleware) *reportController {
	return &reportController{
		ru:             ru,
		rg:             &rg.RouterGroup,
		authMiddleware: auth_middleware,
	}
}
