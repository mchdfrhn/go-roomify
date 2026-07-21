// Package query menyediakan builder query SQL dinamis berparameter aman untuk PostgreSQL.
package query

// Import package database/sql, fmt, dan strings.
import (
	"database/sql" // Interface standar database SQL Go
	"fmt"          // Format string
	"strings"      // Pengolahan string
)

// QSelect merupakan struktur builder query SQL SELECT kompleks.
type QSelect struct {
	DB          *sql.DB  // Pointer koneksi database SQL
	param_index int      // Indeks penomoran placeholder parameter PostgreSQL ($1, $2, dst)
	param_value []any    // Slice daftar nilai argumen parameter query
	qtable      string   // Klausa FROM tabel
	qcolumn     []string // Slice kolom-kolom yang akan di-SELECT
	qjoin       string   // Klausa JOIN (INNER/LEFT)
	qwhere      string   // Klausa WHERE
	qlimit      string   // Klausa LIMIT dan OFFSET
	qoffset     string   // Klausa OFFSET tambahan
	qorder      string   // Klausa ORDER BY
	qgroup      string   // Klausa GROUP BY
}

// Table menentukan tabel sumber data pada klausa FROM.
func (self *QSelect) Table(tbl_name ...string) *QSelect {
	// Menyusun klausa FROM dengan menggabungkan nama tabel jika lebih dari satu
	self.qtable = fmt.Sprintf(" FROM %s", strings.Join(tbl_name, ","))
	// Mengembalikan instance builder
	return self
}

// Column menentukan nama kolom yang diproyeksikan pada klausa SELECT.
func (self *QSelect) Column(col_name ...string) *QSelect {
	// Menambahkan nama kolom ke slice qcolumn
	self.qcolumn = append(self.qcolumn, col_name...)
	// Mengembalikan instance builder
	return self
}

// Where menambahkan kondisi dasar pada klausa WHERE dengan parameter terikat.
func (self *QSelect) Where(condition string, operator string, value any) *QSelect {
	// Menambahkan kata kunci WHERE jika klausa masih kosong
	if self.qwhere == "" {
		self.qwhere = " WHERE"
	}
	// Meng-increment indeks parameter
	self.param_index += 1
	// Menambahkan nilai parameter ke slice argumen
	self.param_value = append(self.param_value, value)
	// Membentuk string ekspresi kondisi berparameter
	self.qwhere += fmt.Sprintf(" %s %s $%d", condition, operator, self.param_index)
	// Mengembalikan instance builder
	return self
}

// WhereColumn menambahkan kondisi antar dua kolom tanpa binding parameter nilai.
func (self *QSelect) WhereColumn(condition string, operator string, value any) *QSelect {
	// Menambahkan kata kunci WHERE jika klausa masih kosong
	if self.qwhere == "" {
		self.qwhere = " WHERE"
	}
	// Membentuk ekspresi perbandingan antar kolom/literal langsung
	self.qwhere += fmt.Sprintf(" %s %s %s", condition, operator, value)
	// Mengembalikan instance builder
	return self
}

// OrWhere menambahkan kondisi logis OR pada klausa WHERE.
func (self *QSelect) OrWhere(condition string, operator string, value any) *QSelect {
	// Memanggil Where dengan prefix OR
	return self.Where(" OR "+condition, operator, value)
}

// AndWhere menambahkan kondisi logis AND pada klausa WHERE.
func (self *QSelect) AndWhere(condition string, operator string, value any) *QSelect {
	// Memanggil Where dengan prefix AND
	return self.Where(" AND "+condition, operator, value)
}

// LeftJoin menyusun klausa LEFT JOIN dengan tabel target dan kondisi relasinya.
func (self *QSelect) LeftJoin(tbl_name string, condition string) *QSelect {
	// Menambahkan pernyataan LEFT JOIN ... ON ...
	self.qjoin += fmt.Sprintf(" LEFT JOIN %s ON %s ", tbl_name, condition)
	// Mengembalikan instance builder
	return self
}

// Join menyusun klausa INNER JOIN dengan tabel target dan kondisi relasinya.
func (self *QSelect) Join(tbl_name string, condition string) *QSelect {
	// Menambahkan pernyataan JOIN ... ON ...
	self.qjoin += fmt.Sprintf(" JOIN %s ON %s ", tbl_name, condition)
	// Mengembalikan instance builder
	return self
}

// Limit menentukan jumlah maksimal baris yang dikembalikan oleh query.
func (self *QSelect) Limit(value int) *QSelect {
	// Meng-increment indeks parameter
	self.param_index += 1
	// Menambahkan nilai limit ke param_value
	self.param_value = append(self.param_value, value)
	// Membentuk klausa LIMIT $x
	self.qlimit += fmt.Sprintf(" LIMIT $%d ", self.param_index)
	// Mengembalikan instance builder
	return self
}

// Offset menentukan offset/baris yang dilewati untuk paginasi.
func (self *QSelect) Offset(value int) *QSelect {
	// Meng-increment indeks parameter
	self.param_index += 1
	// Menambahkan nilai offset ke param_value
	self.param_value = append(self.param_value, value)
	// Membentuk klausa OFFSET $x
	self.qlimit += fmt.Sprintf(" OFFSET $%d ", self.param_index)
	// Mengembalikan instance builder
	return self
}

// GroupBy menentukan pengelompokan baris data berdasarkan nama kolom.
func (self *QSelect) GroupBy(col_name string) *QSelect {
	// Meng-increment indeks parameter
	self.param_index += 1
	// Menambahkan nama kolom pengelompokan
	self.param_value = append(self.param_value, col_name)
	// Membentuk klausa GROUP BY $x
	self.qlimit += fmt.Sprintf(" GROUP BY $%d", self.param_index)
	// Mengembalikan instance builder
	return self
}

// OrderBy menentukan urutan pengurutan data (ASC/DESC).
func (self *QSelect) OrderBy(value string, order string) *QSelect {
	// Meng-increment indeks parameter
	self.param_index += 1
	// Menambahkan nama kolom pengurutan
	self.param_value = append(self.param_value, value)
	// Membentuk klausa ORDER BY $x ASC/DESC
	self.qlimit += fmt.Sprintf(" ORDER BY $%d %s", self.param_index, order)
	// Mengembalikan instance builder
	return self
}

// GetQuery menyusun dan menggabungkan seluruh komponen klausa menjadi perintah SQL SELECT utuh.
func (self *QSelect) GetQuery() string {
	// Memformat query SELECT lengkap
	query := fmt.Sprintf(" SELECT %s %s %s %s %s %s %s %s",
		strings.Join(self.qcolumn, ","),
		self.qtable,
		self.qjoin,
		self.qwhere,
		self.qgroup,
		self.qorder,
		self.qlimit,
		self.qoffset)
	// Mengembalikan string query SQL
	return query
}

// Run mengeksekusi query SELECT yang menghasilkan banyak baris data (*sql.Rows).
func (self *QSelect) Run() (*sql.Rows, error) {
	// Menjalankan Query pada database dengan menyertakan argumen param_value
	rows, err := self.DB.Query(self.GetQuery(), self.param_value...)
	// Mengembalikan cursor rows dan error
	return rows, err
}

// RunRow mengeksekusi query SELECT yang menghasilkan tepat satu baris data (*sql.Row).
func (self *QSelect) RunRow() *sql.Row {
	// Menjalankan QueryRow pada database
	row := self.DB.QueryRow(self.GetQuery(), self.param_value...)
	// Mengembalikan pointer single row
	return row
}
