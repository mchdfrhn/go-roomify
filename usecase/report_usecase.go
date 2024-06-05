package usecase

import (
	"go-roomify/repository"
)

type ReportUsecase interface {
	//DownloadReportByYear( startYear string, endYear string ) ( *sql.Rows, int, error )
}

type reportUsecase struct {
	reservationRepo repository.ReservationRepository
}

// func (ru *reportUsecase) DownloadReportByYear(startYear string, endYear string) (*sql.Rows, int, error) {
// 	if startYear == "" || endYear == "" {
// 		return nil, http.StatusBadRequest, fmt.Errorf("Query parameter 's' or 'e' must not be empty")
// 	}

// 	_, err := time.Parse("2006-1-2", startYear)
// 	if err != nil {
// 		return nil, http.StatusBadRequest, fmt.Errorf("format query s must be YYYY-MM-DD")
// 	}

// 	_, err = time.Parse("2006-1-2", endYear)
// 	if err != nil {
// 		return nil, http.StatusBadRequest, fmt.Errorf("format query e be YYYY-MM-DD")
// 	}

// 	rows, err := ru.reservationRepo.GetReservationByYear( startYear, endYear )
// 	if err != nil {
// 		return nil, http.StatusInternalServerError, err
// 	}

// 	return rows, http.StatusOK, nil

// }

func NewReportUsecase(reservationRepo repository.ReservationRepository) ReportUsecase {
	return &reportUsecase{
		reservationRepo: reservationRepo,
	}
}
