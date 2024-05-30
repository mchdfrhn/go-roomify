package usecase

import (
	"database/sql"
	"fmt"
	"time"
	"go-roomify/repository"
	"go-roomify/utils"
	"net/http"
)

type ReportUsecase interface{
	DownloadReportByYear( startYear, endYear, startMonth, endMonth, startDay, endDay string ) ( *sql.Rows, int, error )
}

type reportUsecase struct{
	reservationRepo repository.ReservationRepository
}

func ( ru *reportUsecase ) DownloadReportByYear( startYear, endYear, startMonth, endMonth, startDay, endDay string ) ( *sql.Rows, int, error ){

	if startYear == "" && endYear != "" || 
		startMonth == "" && endMonth != "" || 
		startDay == "" && endMonth != "" {

		return nil, http.StatusBadRequest, fmt.Errorf("cant set ending without start")
	}

	startYear, endYear = utils.ValidateYear( startYear, endYear )
	startMonth, endMonth = utils.ValidateMonth( startMonth, endMonth )
	startDay, endDay = utils.ValidateDay( startDay, endDay )

	selectedStart := startYear + "-" + startMonth + "-" + startDay
	selectedEnd := endYear + "-" + endMonth + "-" + endDay

	var rows *sql.Rows

	_, err := time.Parse("0000-00-00", selectedEnd)
	if err != nil {
		rows, err = ru.reservationRepo.GetReservationByYear( selectedStart )
		if err != nil {
			return nil, http.StatusInternalServerError, nil
		}

		return rows, http.StatusOK, nil
	}

	rows, err = ru.reservationRepo.GetReservationByYearBetween( selectedStart, selectedEnd )
	if err != nil {
		return nil, http.StatusInternalServerError, nil
	}

	return rows, http.StatusOK, nil
}

func NewReportUsecase( reservationRepo repository.ReservationRepository ) ReportUsecase{
	return &reportUsecase{
		reservationRepo: reservationRepo,
	}
}