// Package utils menyediakan sekumpulan fungsi pembantu (utility) untuk aplikasi go-roomify.
package utils

// Import DTO, math, os, strconv, dan godotenv.
import (
	"go-roomify/model/dto" // Import struct DTO paginasi
	"math"                 // Import fungsi matematika (Ceil)
	"os"                   // Pembacaan variabel lingkungan
	"strconv"              // Konversi string ke angka

	"github.com/joho/godotenv" // Driver pemicu pembaca file .env
)

// getPaginationParams menghitung parameter halaman (Page, Take, Skip) untuk query database.
func getPaginationParams(params dto.PaginationParam) dto.PaginationQuery {
	// Inisialisasi variabel batas halaman
	var page int
	var take int
	var skip int

	// Memeriksa jika nomor halaman lebih dari 1
	if params.Page > 1 {
		page = params.Page // Menggunakan nomor halaman yang diminta
	} else {
		page = 1 // Default ke halaman pertama jika <= 1
	}

	// Memeriksa apakah limit data bernilai 0 (tidak ditentukan)
	if params.Limit == 0 {
		// Memuat variabel .env
		err := godotenv.Load(".env")
		if err != nil {
			// Mengembalikan struct kosong jika file .env gagal dimuat
			return dto.PaginationQuery{}
		}
		// Mengambil nilai default rows per page dari environment
		n, _ := strconv.Atoi(os.Getenv("DEFAULT_ROWS_PER_PAGE"))
		take = n
	} else {
		take = params.Limit // Menggunakan limit yang ditentukan
	}

	// Menghitung jumlah record yang akan di-skip (OFFSET = (Page - 1) * Take)
	skip = (page - 1) * take

	// Mengembalikan struct PaginationQuery hasil perhitungan
	return dto.PaginationQuery{
		Page: page,
		Take: take,
		Skip: skip,
	}
}

// Paginate membuat metadata DTO Paging lengkap berdasarkan total baris dan limit per halaman.
func Paginate(page, limit, totalRows int) dto.Paging {
	// Mengembalikan objek DTO Paging
	return dto.Paging{
		Page:        page,                                                   // Halaman saat ini
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(limit))), // Menghitung total halaman dengan pembulatan ke atas
		TotalRows:   totalRows,                                              // Total keseluruhan baris data
		RowsPerPage: limit,                                                  // Jumlah baris per halaman
	}
}