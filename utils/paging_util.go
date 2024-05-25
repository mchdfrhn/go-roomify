package utils

import (
	"go-roomify/model/dto"
	"math"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func getPaginationParams(params dto.PaginationParam) dto.PaginationQuery {
	var page int
	var take int
	var skip int

	if params.Page > 1 {
		page = params.Page
	}else{
		page = 1
	}

	if params.Limit == 0 {
		err := godotenv.Load(".env")
		if err != nil {
			return dto.PaginationQuery{}
		}
		n, _ := strconv.Atoi(os.Getenv("DEFAULT_ROWS_PER_PAGE"))
		take = n
	}else {
		take = params.Limit
	}

	skip = (page - 1) * take

	return dto.PaginationQuery{
		Page : page,
		Take : take,
		Skip : skip,
	}
}

func Paginate(page,limit,totalRows int) dto.Paging {
	return dto.Paging {
		Page : page,
		TotalPages: int(math.Ceil(float64(totalRows)/float64(limit))),
		TotalRows : totalRows,
		RowsPerPage : limit,
	}
}