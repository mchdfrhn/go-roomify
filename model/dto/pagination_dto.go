// Package dto berisi struktur Data Transfer Object untuk permintaan dan tanggapan API.
package dto

// QueryParams menyimpan parameter pencarian (search query), pengurutan kolom (order), dan arah urutan (sort).
type QueryParams struct {
	Query string // Kata kunci pencarian data
	Order string // Nama kolom yang diurutkan
	Sort  string // Arah pengurutan (asc atau desc)
}

// isSortValid mengecek apakah nilai opsi pengurutan sort bernilai 'asc' atau 'desc'.
func (qp *QueryParams) isSortValid() bool {
	// Memeriksa jika sort bernilai "asc" atau "desc"
	return qp.Sort == "asc" || qp.Sort == "desc"
}

// PaginationParam menyimpan parameter halaman, offset baris data, dan limit jumlah data.
type PaginationParam struct {
	Page   int // Nomor halaman aktif saat ini
	Offset int // Jumlah baris yang dilewati dalam query SQL
	Limit  int // Jumlah maksimum baris yang diambil
}

// PaginationQuery menyimpan struktur query pagination alternatif.
type PaginationQuery struct {
	Page int // Nomor halaman
	Take int // Jumlah data yang diambil (limit)
	Skip int // Jumlah data yang dilewati (offset)
}

// Paging menyimpan informasi meta-data halaman untuk dikembalikan pada respon JSON.
type Paging struct {
	Page        int // Halaman saat ini
	RowsPerPage int // Jumlah baris per halaman
	TotalRows   int // Total seluruh baris data di database
	TotalPages  int // Total jumlah halaman yang tersedia
}

// RequestQueryParam mengombinasikan QueryParams dan PaginationParam untuk penanganan query string request.
type RequestQueryParam struct {
	QueryParams     // Embed QueryParams
	PaginationParam // Embed PaginationParam
}