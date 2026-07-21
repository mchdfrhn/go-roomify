// Package utils menyediakan sekumpulan fungsi pembantu (utility) untuk aplikasi go-roomify.
package utils

// Import fmt, dto, strconv, dan gin.
import (
	"fmt"                  // Format error dan string
	"go-roomify/model/dto" // DTO request dan pagination
	"strconv"              // Konversi string ke int

	"github.com/gin-gonic/gin" // Context HTTP framework Gin
)

// validateRequestQueryParams memvalidasi dan memposting parameter query URL (page, limit, order, sort).
func validateRequestQueryParams(c *gin.Context) (dto.RequestQueryParam, error) {

	// Validasi dan konversi query parameter "page", default nilai "1"
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	// Memeriksa jika terjadi kesalahan konversi atau nomor page kurang dari sama dengan 0
	if err != nil || page <= 0 {
		// Mengembalikan struct kosong dan error "Invalid Page Number"
		return dto.RequestQueryParam{}, fmt.Errorf("Invalid Page NUmber")
	}

	// Validasi dan konversi query parameter "limit", default nilai "5"
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	// Memeriksa jika terjadi kesalahan konversi atau limit kurang dari sama dengan 0
	if err != nil || limit <= 0 {
		// Mengembalikan struct kosong dan error "Invalid Limit Number"
		return dto.RequestQueryParam{}, fmt.Errorf("Invalid Limit NUmber")
	}

	// Mengambil query parameter "order", default "id"
	order := c.DefaultQuery("order", "id")
	// Mengambil query parameter "sort", default "asc"
	sort := c.DefaultQuery("sort", "asc")

	// Mengembalikan struct RequestQueryParam yang terisi beserta nil error
	return dto.RequestQueryParam{
		QueryParams: dto.QueryParams{
			Order: order,
			Sort:  sort,
		},
		PaginationParam: dto.PaginationParam{
			Page:  page,
			Limit: limit,
		},
	}, nil
}