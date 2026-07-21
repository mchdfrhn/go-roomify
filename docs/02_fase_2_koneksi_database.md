# 🔌 Fase 2: Menghubungkan ke Database (PostgreSQL)

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [⬅️ Fase 1](01_fase_1_hello_world_server.md) | [Lanjut ke Fase 3 ➡️](03_fase_3_crud_pertama_monolit.md)

---

## 💡 Mindset Fase 2

Setelah web server berjalan, aplikasi biasanya butuh penyimpan data (Database).  
Di tahap ini, kita akan:
1. Meng-install driver PostgreSQL (`github.com/lib/pq`).
2. Menambahkan fungsi tes koneksi database langsung di dalam file proyek kita.

---

## 1. Install Driver PostgreSQL & Godotenv

```bash
# Driver PostgreSQL untuk Go database/sql
go get -u github.com/lib/pq

# Manajer Environment Variable
go get -u github.com/joho/godotenv
```

---

## 2. Buat File Environment (`.env`)

Buat file `.env` di root folder aplikasi:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=roomify_db
DB_DRIVER=postgres
API_PORT=:8085
```

---

## 3. Tambahkan Koneksi DB ke `main.go`

Ubah `main.go` Anda menjadi seperti ini untuk menguji apakah Go dapat berkomunikasi dengan PostgreSQL:

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Driver PostgreSQL
)

func main() {
	// 1. Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan env sistem")
	}

	// 2. Format DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	// 3. Buka koneksi database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi ke DB: ", err)
	}
	defer db.Close()

	// 4. Ping database untuk memastikan koneksi hidup
	if err := db.Ping(); err != nil {
		log.Fatal("Database Ping Error: ", err)
	}

	fmt.Println("✅ Berhasil terhubung ke database PostgreSQL!")

	// 5. Inisialisasi HTTP Router
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "OK",
			"database": "Connected",
		})
	})

	r.Run(os.Getenv("API_PORT"))
}
```

---

## 4. Coba Jalankan

```bash
go run main.go
```

**Output Terminal:**
```text
✅ Berhasil terhubung ke database PostgreSQL!
[GIN-debug] Listening and serving HTTP on :8085
```

✅ **Hebat!** Aplikasi Anda sekarang sudah terhubung dengan aman ke Database PostgreSQL. Tahap berikutnya adalah **Fase 3: Membuat Fitur CRUD Pertama**.

---

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [Lanjut ke Fase 3 ➡️](03_fase_3_crud_pertama_monolit.md)
