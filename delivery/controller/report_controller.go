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

	ctx.Header("Content-type", "text/csv")
	ctx.Header("Content-Disposition", "attachment; filename=\"report_start="+startYear+"_end="+endYear+".csv\"")

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