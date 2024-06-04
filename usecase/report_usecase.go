package usecase

import (
	"database/sql"
	"fmt"
	"time"
	"go-roomify/repository"
	"net/http"
)

type ReportUsecase interface{
	DownloadReportByYear( startYear string, endYear string ) ( *sql.Rows, int, error )
}

type reportUsecase struct{
	reportRepo repository.ReportRepository
}

func (ru *reportUsecase) DownloadReportByYear(startYear string, endYear string) (*sql.Rows, int, error) {
	if startYear == "" || endYear == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("query parameter 's' or 'e' must not be empty")
	}

	_, err := time.Parse("2006-1-2", startYear)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("format query s must be YYYY-MM-DD")
	}

	_, err = time.Parse("2006-1-2", endYear)
	if err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("format query e be YYYY-MM-DD")
	}

	rows, err := ru.reportRepo.GetReportByYear( startYear, endYear )
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return rows, http.StatusOK, nil

}


func NewReportUsecase( reportRepo repository.ReportRepository ) ReportUsecase{
	return &reportUsecase{
		reportRepo: reportRepo,
	}
}