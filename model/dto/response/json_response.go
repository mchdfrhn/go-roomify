// Package response menyediakan utility helper untuk mengirimkan standar response JSON HTTP.
package response

// Import dependency Gin framework, HTTP status code, dan DTO paging.
import (
	"go-roomify/model/dto" // Import DTO pagination
	"net/http"            // Import paket HTTP untuk konstanta status code

	"github.com/gin-gonic/gin" // Framework Gin web engine
)

// SendSingleResponseCreated mengembalikan JSON response dengan status HTTP 201 Created (atau 200 OK dengan kode status data 201).
func SendSingleResponseCreated(ctx *gin.Context, data any, descriptionMsg string) {
	// Mengirimkan payload JSON SingleResponse dengan status 201
	ctx.JSON(http.StatusOK, &SingleResponse{
		Status: Status{
			Code:        http.StatusCreated, // Status code 201 Created
			Description: descriptionMsg,      // Pesan deskripsi sukses
		},
		Data: data, // Data payload tunggal yang dihasilkan
	})
}

// SendSingleResponseData mengembalikan JSON response standar 200 OK beserta data payload.
func SendSingleResponseData(ctx *gin.Context, data any, descriptionMsg string) {
	// Mengirimkan payload JSON SingleResponse dengan status 200 OK
	ctx.JSON(http.StatusOK, &SingleResponse{
		Status: Status{
			Code:        http.StatusOK,  // Status code 200 OK
			Description: descriptionMsg, // Pesan deskripsi sukses
		},
		Data: data, // Data payload tunggal
	})
}

// SendSinglePageResponse mengembalikan JSON response berpaginasi (paged response) dengan meta data Paging.
func SendSinglePageResponse(ctx *gin.Context, data []any, descriptionMsg string, paging dto.Paging) {
	// Mengirimkan payload JSON PagedResponse lengkap dengan array data dan metadata paging
	ctx.JSON(http.StatusOK, &PagedResponse{
		Status: Status{
			Code:        http.StatusOK,  // Status code 200 OK
			Description: descriptionMsg, // Pesan deskripsi sukses
		},
		Data:   data,   // Slice list data hasil query
		Paging: paging, // Metadata paginasi (total rows, pages, dsb)
	})
}

// SendSingleResponse mengembalikan JSON response standar hanya berisi status dan pesan deskripsi tanpa data payload.
func SendSingleResponse(ctx *gin.Context, descriptionMsg string) {
	// Mengirimkan status deskripsi JSON
	ctx.JSON(http.StatusOK, &Status{
		Code:        http.StatusOK,  // Status code 200 OK
		Description: descriptionMsg, // Pesan deskripsi
	})
}

// SendSingleResponseError menghentikan eksekusi handler (Abort) dan mengembalikan JSON error sesuai HTTP status code.
func SendSingleResponseError(ctx *gin.Context, code int, errorMessage string) {
	// Membatalkan pipeline handler middleware Gin dan mengirimkan JSON error status
	ctx.AbortWithStatusJSON(code, &Status{
		Code:        code,         // Status code error (400, 401, 403, 404, 500, dsb)
		Description: errorMessage, // Pesan rincian error
	})
}
