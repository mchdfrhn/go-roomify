// Package query menyediakan builder query SQL dinamis berparameter aman untuk PostgreSQL.
package query

// Import package database/sql, fmt, dan strings.
import (
	"database/sql" // Interface standar database SQL Go
	"fmt"          // Format string
	"strings"      // Pengolahan string
)

// QUpdate merupakan struktur builder query SQL UPDATE.
type QUpdate struct {
	DB          *sql.DB  // Pointer koneksi database SQL
	param_index int      // Indeks penomoran placeholder parameter PostgreSQL ($1, $2, dst)
	param_value []any    // Slice daftar nilai argumen parameter yang diperbarui
	qtable      string   // Nama tabel target pembaruan
	qset        []string // Slice pasangan penetapan kolom ("column = $1")
	qwhere      string   // Klausa WHERE
}

// Table menentukan nama tabel yang data barisnya akan diperbarui.
func (self *QUpdate) Table(tbl_name string) *QUpdate {
	// Menyimpan nama tabel target
	self.qtable = tbl_name
	// Mengembalikan instance builder untuk method chaining
	return self
}

// Set menetapkan pasangan kolom dan nilai yang akan diperbarui pada klausa SET.
func (self *QUpdate) Set(col_name string, value any) *QUpdate {
	// Meng-increment indeks parameter
	self.param_index += 1
	// Menambahkan nilai baru ke slice param_value
	self.param_value = append(self.param_value, value)
	// Membentuk string pasangan "nama_kolom = $x"
	self.qset = append(self.qset, fmt.Sprintf("%s = $%d", col_name, self.param_index))
	// Mengembalikan instance builder
	return self
}

// Where menambahkan kondisi kriteria data yang diperbarui pada klausa WHERE.
func (self *QUpdate) Where(condition string, operator string, value any) *QUpdate {
	// Menambahkan kata kunci WHERE jika klausa masih kosong
	if self.qwhere == "" {
		self.qwhere = " WHERE"
	}
	// Meng-increment indeks parameter
	self.param_index += 1
	// Menambahkan nilai parameter ke slice argumen
	self.param_value = append(self.param_value, value)
	// Membentuk klausa WHERE berparameter
	self.qwhere += fmt.Sprintf(" %s %s $%d", condition, operator, self.param_index)
	// Mengembalikan instance builder
	return self
}

// OrWhere menambahkan kondisi logis OR pada klausa WHERE UPDATE.
func (self *QUpdate) OrWhere(condition string, operator string, value any) *QUpdate {
	// Memanggil Where dengan prefix OR
	return self.Where(" OR "+condition, operator, value)
}

// AndWhere menambahkan kondisi logis AND pada klausa WHERE UPDATE.
func (self *QUpdate) AndWhere(condition string, operator string, value any) *QUpdate {
	// Memanggil Where dengan prefix AND
	return self.Where(" AND "+condition, operator, value)
}

// GetQuery menyusun dan mengembalikan string statement SQL UPDATE lengkap.
func (self *QUpdate) GetQuery() string {
	// Memformat query UPDATE SET WHERE
	query := fmt.Sprintf(" UPDATE %s SET %s %s",
		self.qtable,
		strings.Join(self.qset, ","),
		self.qwhere)
	// Mengembalikan string query SQL
	return query
}

// Run mengeksekusi query UPDATE pada database menggunakan parameter yang aman.
func (self *QUpdate) Run() (sql.Result, error) {
	// Menjalankan Exec untuk memperbarui baris data
	result, err := self.DB.Exec(self.GetQuery(), self.param_value...)
	// Mengembalikan hasil eksekusi (RowsAffected) dan error jika ada
	return result, err
}
