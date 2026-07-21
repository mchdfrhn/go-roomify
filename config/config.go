// Package config menangani pembacaan dan pengelolaan konfigurasi aplikasi dari environment variable (.env).
package config

// Import package standar Go dan library eksternal godotenv untuk membaca file .env.
import (
	"errors" // Package untuk membuat error sederhana
	"fmt"    // Package untuk format string dan error
	"os"     // Package untuk mengakses variabel lingkungan OS (environment variables)
	"strconv" // Package untuk konversi tipe data string ke angka
	"time"   // Package untuk penanganan durasi waktu

	"github.com/joho/godotenv" // Library pihak ketiga untuk membaca file .env
)

// DbConfig menyimpan data konfigurasi koneksi ke database PostgreSQL.
type DbConfig struct {
	Host     string // Alamat host server database (contoh: localhost / IP)
	Port     string // Nomor port server database (contoh: 5432)
	Name     string // Nama database yang digunakan
	User     string // Username autentikasi database
	Password string // Kata sandi autentikasi database
	Driver   string // Driver database yang digunakan (contoh: postgres)
}

// Config merupakan struktur utama yang menggabungkan seluruh konfigurasi aplikasi.
type Config struct {
	DbConfig    // Embed struktur DbConfig untuk konfigurasi database
	TokenConfig // Embed struktur TokenConfig untuk konfigurasi token JWT
}

// TokenConfig menyimpan konfigurasi terkait autentikasi JWT (JSON Web Token).
type TokenConfig struct {
	IssuerName      string        // Nama penerbit token (issuer)
	JwtSignatureKey []byte        // Kunci rahasia (secret key) untuk penandatanganan JWT
	JwtLifeTime     time.Duration // Durasi masa berlaku token JWT
}

// ReadConfig membaca konfigurasi dari file .env dan memasukkannya ke dalam struct Config.
func (c *Config) ReadConfig() error {
	// Memuat variabel dari file .env ke environment aplikasi.
	err := godotenv.Load()
	// Memeriksa apakah terjadi error saat membaca file .env.
	if err != nil {
		// Mengembalikan error jika file .env gagal dimuat.
		return fmt.Errorf("Error Load .env file")
	}

	// Mengisi struktur DbConfig dari variabel environment yang dibaca.
	c.DbConfig = DbConfig{
		Host:     os.Getenv("DB_HOST"),   // Mengambil nilai DB_HOST
		Port:     os.Getenv("DB_PORT"),   // Mengambil nilai DB_PORT
		Name:     os.Getenv("DB_NAME"),   // Mengambil nilai DB_NAME
		User:     os.Getenv("DB_USER"),   // Mengambil nilai DB_USER
		Password: os.Getenv("DB_PASS"),   // Mengambil nilai DB_PASS
		Driver:   os.Getenv("DB_DRIVER"), // Mengambil nilai DB_DRIVER
	}

	// Mengonversi string TOKEN_LIFE_TIME menjadi integer (jam).
	token_lifetime, err := strconv.Atoi(os.Getenv("TOKEN_LIFE_TIME"))
	// Memeriksa apakah terjadi error konversi string ke int.
	if err != nil {
		// Mengembalikan error jika parsing durasi token gagal.
		return errors.New("Fail Parse Token Life Time")
	}

	// Mengisi struktur TokenConfig dengan data dari environment.
	c.TokenConfig = TokenConfig{
		IssuerName:      os.Getenv("ISSUER_NAME"),                   // Mengambil nama issuer token
		JwtSignatureKey: []byte(os.Getenv("SIGNATURE")),             // Mengambil secret key JWT sebagai byte array
		JwtLifeTime:     time.Duration(token_lifetime) * time.Hour, // Menghitung durasi token dalam satuan jam
	}

	// Memeriksa apakah ada variabel environment database wajib yang kosong.
	if c.DbConfig.Host == "" || c.DbConfig.Port == "" || c.DbConfig.Name == "" || c.DbConfig.User == "" || c.DbConfig.Password == "" || c.DbConfig.Driver == "" {
		// Mengembalikan error jika ada variabel environment yang belum terisi.
		return fmt.Errorf("Missing ENV")
	}
	// Mengembalikan nil jika konfigurasi berhasil dibaca tanpa error.
	return nil
}

// NewConfig membuat instance Config baru dan langsung membaca file konfigurasi.
func NewConfig() (*Config, error) {
	// Membuat alokasi memory baru untuk pointer struct Config.
	cfg := &Config{}
	// Memanggil metode ReadConfig untuk membaca isi environment variable.
	err := cfg.ReadConfig()
	// Memeriksa apakah terjadi error saat proses ReadConfig.
	if err != nil {
		// Mengembalikan nil dan error jika konfigurasi gagal dibaca.
		return nil, err
	}
	// Mengembalikan pointer objek Config dan nil error jika berhasil.
	return cfg, nil
}
