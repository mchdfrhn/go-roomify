package controller

import (
	"go-roomify/model/dto/response"
	"go-roomify/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/sqltocsv"
)

type reportController struct{
	ru usecase.ReportUsecase
	rg *gin.RouterGroup
}

func( rc *reportController ) downloadReportByYearHandler(ctx *gin.Context){

	ctx.Header("Content-type", "text/csv")
	ctx.Header("Content-Disposition", "attachment; filename=\"report.csv\"")

	startYear := ctx.Query("sy")
	endYear := ctx.Query("ey")

	startMonth := ctx.Query("sm")
	endMonth := ctx.Query("em")

	startDay := ctx.Query("sd")
	endDay := ctx.Query("ed")

	rows, code, err := rc.ru.DownloadReportByYear( 
		startYear, 
		endYear,
		startMonth,
		endMonth,
		startDay,
		endDay, 
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

func( rc *reportController ) Route(){
	group := rc.rg.Group("/report")
	group.GET("/", rc.downloadReportByYearHandler)
}

func NewReportController( ru usecase.ReportUsecase, rg *gin.Engine ) *reportController{
	return &reportController{
		ru: ru,
		rg: &rg.RouterGroup,
	}
}