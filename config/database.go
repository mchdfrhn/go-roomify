// Package config menangani konfigurasi dan pembentukan koneksi ke PostgreSQL database.
package config

// Import package database/sql untuk interaksi database standar Go dan driver postgres.
import (
	"database/sql" // Package standar Go untuk SQL database interface
	"fmt"          // Package untuk pembentukan string terformat

	_ "github.com/lib/pq" // Driver PostgreSQL (side-effect import untuk registrasi driver "postgres")
)

// DbConnection merupakan interface pengabstraksi untuk mendapatkan pointer ke *sql.DB.
type DbConnection interface {
	Conn() *sql.DB // Method untuk mengembalikan koneksi *sql.DB yang sudah aktif
}

// dbConnection merupakan implementasi konkret dari interface DbConnection.
type dbConnection struct {
	db  *sql.DB // Object koneksi database SQL
	cfg *Config // Reference ke konfigurasi aplikasi
}

// initDb menginisialisasi koneksi ke PostgreSQL menggunakan konfigurasi DSN.
func (d *dbConnection) initDb() error {
	// Membentuk string DSN (Data Source Name) sesuai spesifikasi lib/pq.
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.cfg.DbConfig.Host,     // Alamat host database
		d.cfg.DbConfig.Port,     // Port database
		d.cfg.DbConfig.User,     // Username database
		d.cfg.DbConfig.Password, // Password database
		d.cfg.DbConfig.Name,     // Nama database
	)
	// Membuka handle database dengan driver dan dsn yang disiapkan.
	db, err := sql.Open(d.cfg.DbConfig.Driver, dsn)
	// Memeriksa jika terjadi kegagalan pembukaan handle database.
	if err != nil {
		// Mengembalikan error terformat jika sql.Open gagal.
		return fmt.Errorf("Error Create Connection : %v", err)
	}
	// Melakukan test ping ke server database untuk memastikan koneksi benar-benar terhubung.
	checkConn := db.Ping()
	// Memeriksa hasil ping koneksi.
	if checkConn != nil {
		// Mengembalikan error jika server database tidak merespons ping.
		return fmt.Errorf("Error Ping !! : %v", checkConn.Error())
	}
	// Menyimpan instance koneksi *sql.DB yang berhasil dibuat ke dalam field db.
	d.db = db
	// Mengembalikan nil tanda inisialisasi berhasil.
	return nil
}

// Conn mengembalikan instance koneksi *sql.DB yang aktif.
func (d *dbConnection) Conn() *sql.DB {
	// Mengembalikan pointer d.db.
	return d.db
}

// NewDbConnection menginisialisasi dbConnection baru dengan konfigurasi yang diberikan.
func NewDbConnection(config *Config) (DbConnection, error) {
	// Membuat instance dbConnection baru dengan menyimpan referensi config.
	conn := &dbConnection{
		cfg: config,
	}

	// Memanggil metode initDb untuk membuka koneksi dan verifikasi ping.
	err := conn.initDb()
	// Memeriksa jika terdapat error saat initDb.
	if err != nil {
		// Mengembalikan nil dan error jika koneksi gagal dibentuk.
		return nil, err
	}
	// Mengembalikan interface DbConnection dan nil error jika sukses.
	return conn, nil
}
