// Package query menyediakan builder query SQL dinamis berparameter aman untuk PostgreSQL.
package query

// Import package database/sql, fmt, dan strings.
import (
	"database/sql" // Interface standar database SQL Go
	"fmt"          // Format string
	"strings"      // Pengolahan string dan gabungan slice
)

// QInsert merupakan struktur builder query SQL INSERT INTO.
type QInsert struct {
	DB          *sql.DB  // Pointer koneksi database SQL
	param_index int      // Indeks penomoran placeholder parameter PostgreSQL ($1, $2, dst)
	param_value []any    // Slice daftar nilai argumen yang disisipkan
	qtable      string   // Nama tabel tujuan penyisipan
	qcolumn     []string // Slice nama-nama kolom target
	qvalues     []string // Slice ekspresi himpunan nilai ($1, $2, ...)
}

// Table menentukan nama tabel target penyisipan data.
func (self *QInsert) Table(tbl_name string) *QInsert {
	// Menyimpan nama tabel target
	self.qtable = tbl_name
	// Mengembalikan instance builder untuk method chaining
	return self
}

// Column menentukan daftar kolom tempat nilai akan dimasukkan.
func (self *QInsert) Column(col_name ...string) *QInsert {
	// Menambahkan daftar nama kolom ke slice qcolumn
	self.qcolumn = append(self.qcolumn, col_name...)
	// Mengembalikan instance builder
	return self
}

// Values menentukan daftar nilai yang dimasukkan dengan pembuatan placeholder parameter terikat.
func (self *QInsert) Values(col_value ...any) *QInsert {
	// Slice sementara untuk menyimpan string placeholder baris ini ($1, $2, ...)
	var qvalues_str []string
	// Iterasi setiap nilai yang dikirimkan
	for _, value := range col_value {
		// Meng-increment indeks parameter
		self.param_index += 1
		// Menambahkan nilai ke slice param_value
		self.param_value = append(self.param_value, value)
		// Membentuk string placeholder (contoh: $1)
		qvalues_str = append(qvalues_str, fmt.Sprintf("$%d", self.param_index))
	}
	// Menggabungkan string placeholder dalam tanda kurung "( $1, $2 )"
	self.qvalues = append(self.qvalues, "("+strings.Join(qvalues_str, ",")+")")
	// Mengembalikan instance builder
	return self
}

// GetQuery menyusun dan mengembalikan string statement SQL INSERT INTO dasar.
func (self *QInsert) GetQuery() string {
	// Memformat query SQL INSERT INTO
	query := fmt.Sprintf(" INSERT INTO %s (%s) VALUES %s",
		self.qtable,
		strings.Join(self.qcolumn, ","),
		strings.Join(self.qvalues, ","))
	// Mengembalikan string query
	return query
}

// GetQueryReturn menyusun string statement SQL INSERT INTO dengan klausa RETURNING.
func (self *QInsert) GetQueryReturn(col_return ...string) string {
	// Memformat query SQL INSERT INTO beserta klausa RETURNING
	query := fmt.Sprintf(" INSERT INTO %s (%s) VALUES %s RETURNING %s",
		self.qtable,
		strings.Join(self.qcolumn, ","),
		strings.Join(self.qvalues, ","),
		strings.Join(col_return, ","))
	// Mengembalikan string query ber-RETURNING
	return query
}

// RunReturn mengeksekusi query INSERT dan mengembalikan single row (sql.Row) dari klausa RETURNING.
func (self *QInsert) RunReturn(col_return ...string) *sql.Row {
	// Menjalankan QueryRow dengan menyuntikkan parameter nilai
	row := self.DB.QueryRow(self.GetQueryReturn(col_return...), self.param_value...)
	// Mengembalikan pointer sql.Row
	return row
}

// Run mengeksekusi query INSERT biasa pada database.
func (self *QInsert) Run() (sql.Result, error) {
	// Menjalankan Exec untuk menyisipkan data
	result, err := self.DB.Exec(self.GetQuery(), self.param_value...)
	// Mengembalikan hasil eksekusi dan error jika ada
	return result, err
}
