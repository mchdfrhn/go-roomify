// Package query menyediakan builder query SQL dinamis berparameter aman untuk PostgreSQL.
package query

// Import package database/sql dan fmt.
import (
	"database/sql" // Interface standar database SQL Go
	"fmt"          // Format string
)

// QDelete merupakan struktur builder query SQL DELETE.
type QDelete struct {
	DB          *sql.DB // Pointer koneksi database SQL
	param_index int     // Indeks penomoran placeholder parameter PostgreSQL ($1, $2, dst)
	param_value []any   // Slice daftar nilai argumen parameter yang di-binding
	qtable      string  // Nama tabel target penghapusan
	qwhere      string  // Klausa kondisi WHERE
}

// Table menentukan nama tabel tempat baris data akan dihapus.
func (self *QDelete) Table(tbl_name string) *QDelete {
	// Menyimpan nama tabel target
	self.qtable = tbl_name
	// Mengembalikan instance builder untuk mendukung method chaining
	return self
}

// Where menambahkan kondisi kriteria penghapusan pada klausa WHERE.
func (self *QDelete) Where(condition string, operator string, value any) *QDelete {
	// Menambahkan kata kunci WHERE jika klausa masih kosong
	if self.qwhere == "" {
		self.qwhere = " WHERE"
	}
	// Meng-increment indeks parameter ($1, $2, ...)
	self.param_index += 1
	// Menambahkan nilai parameter ke slice argumen
	self.param_value = append(self.param_value, value)
	// Membentuk klausa ekspresi WHERE dengan placeholder parameter yang aman
	self.qwhere += fmt.Sprintf(" %s %s $%d", condition, operator, self.param_index)
	// Mengembalikan instance builder
	return self
}

// OrWhere menambahkan kondisi alternatif dengan operator logis OR pada klausa WHERE.
func (self *QDelete) OrWhere(condition string, operator string, value any) *QDelete {
	// Memanggil fungsi Where dengan imbuhan operator " OR "
	return self.Where(" OR "+condition, operator, value)
}

// AndWhere menambahkan kondisi tambahan dengan operator logis AND pada klausa WHERE.
func (self *QDelete) AndWhere(condition string, operator string, value any) *QDelete {
	// Memanggil fungsi Where dengan imbuhan operator " AND "
	return self.Where(" AND "+condition, operator, value)
}

// GetQuery menyusun dan mengembalikan string perintah query DELETE SQL lengkap.
func (self *QDelete) GetQuery() string {
	// Memformat string statement SQL DELETE
	query := fmt.Sprintf(" DELETE FROM %s %s",
		self.qtable,
		self.qwhere)
	// Mengembalikan string query SQL
	return query
}

// Run mengeksekusi query DELETE pada database menggunakan parameter yang aman.
func (self *QDelete) Run() (sql.Result, error) {
	// Menjalankan query Exec dengan variadic parameter
	result, err := self.DB.Exec(self.GetQuery(), self.param_value...)
	// Mengembalikan hasil eksekusi (RowsAffected) dan error jika ada
	return result, err
}
