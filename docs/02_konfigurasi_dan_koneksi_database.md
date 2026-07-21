# ⚙️ Bagian 2: Konfigurasi Aplikasi & Koneksi Database (`config/`)

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [⬅️ Bagian 1](01_inisialisasi_dan_dependensi.md) | [Lanjut ke Bagian 3 ➡️](03_model_domain_dan_dto.md)

---

## 1. File Environment (`.env`)

File `.env` digunakan untuk menyimpan kredensial dan konfigurasi sensitif agar tidak tercampur di kode sumber (*source code*).

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=roomify_db
DB_DRIVER=postgres

API_PORT=:8085
TOKEN_SECRET=SecretKeySangatRahasia123!
TOKEN_EXPIRE=24h
TOKEN_ISSUER=go-roomify
```

---

## 2. Reading Config (`config/config.go`)

Struct ini bertindak sebagai wadah tersentralisasi untuk semua variabel aplikasi yang dibaca menggunakan `godotenv`.

```go
package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DbConfig    DbConfig
	ApiPort     string
	TokenConfig TokenConfig
}

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Driver   string
}

type TokenConfig struct {
	JwtSignature string
	JwtLifeTime  time.Duration
	JwtIssuer    string
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load() // Membaca file .env di root folder

	return &Config{
		DbConfig: DbConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Driver:   os.Getenv("DB_DRIVER"),
		},
		ApiPort: os.Getenv("API_PORT"),
		TokenConfig: TokenConfig{
			JwtSignature: os.Getenv("TOKEN_SECRET"),
			JwtLifeTime:  24 * time.Hour,
			JwtIssuer:    os.Getenv("TOKEN_ISSUER"),
		},
	}, nil
}
```

---

## 3. Database Connection (`config/db_connection.go`)

File ini membuka koneksi ke PostgreSQL menggunakan `database/sql` & `lib/pq`, dan menyediakan metode `Conn()` untuk mendapatkan instance `*sql.DB`.

```go
package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DbConnection interface {
	Conn() *sql.DB
}

type dbConnection struct {
	db *sql.DB
}

func (d *dbConnection) Conn() *sql.DB {
	return d.db
}

func NewDbConnection(cfg *Config) (DbConnection, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbConfig.Host, cfg.DbConfig.Port, cfg.DbConfig.User, cfg.DbConfig.Password, cfg.DbConfig.Name,
	)

	db, err := sql.Open(cfg.DbConfig.Driver, dsn)
	if err != nil {
		return nil, err
	}

	// Tes koneksi database
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &dbConnection{db: db}, nil
}
```

---

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md)
