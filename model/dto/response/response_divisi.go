// Package response menyediakan helper pengiriman respon JSON spesifik.
package response

// Import model dan paket gin framework.
import (
	"go-roomify/model" // Import model entitas

	"github.com/gin-gonic/gin" // Import Gin Web Framework
)

// SendDivisionsResponse mengembalikan response JSON berupa slice daftar divisi.
func SendDivisionsResponse(ctx *gin.Context, statusCode int, divisions []model.Division) {
	// Mengirimkan payload JSON dengan key 'divisions' dan kode status HTTP tertentu
	ctx.JSON(statusCode, gin.H{"divisions": divisions})
}

// SendDivisionResponse mengembalikan response JSON berupa data objek tunggal divisi.
func SendDivisionResponse(ctx *gin.Context, statusCode int, division model.Division) {
	// Mengirimkan payload JSON dengan key 'division' dan kode status HTTP tertentu
	ctx.JSON(statusCode, gin.H{"division": division})
}
