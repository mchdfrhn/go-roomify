# 🛠️ Fase 3: Membuat Fitur CRUD Pertama (Monolitik / 1 File)

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [⬅️ Fase 2](02_fase_2_koneksi_database.md) | [Lanjut ke Fase 4 ➡️](04_fase_4_refactoring_clean_architecture.md)

---

## 💡 Mindset Fase 3

Sebelum membuat arsitektur terpisah yang rumit, seorang pengembang biasanya **membuat 1 fitur CRUD sampai berfungsi lebih dulu** di file utama.

Ini membantu kita memahami alur lengkap data:
`HTTP JSON Request ➔ Validasi ➔ Query SQL ➔ Scanning DB ➔ HTTP JSON Response`.

---

## 1. Siapkan Tabel di Database PostgreSQL

Jalankan perintah SQL ini di PostgreSQL (pgAdmin / DBeaver / psql):

```sql
CREATE TABLE mst_room (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    capacity INT NOT NULL,
    is_available BOOLEAN DEFAULT true
);
```

---

## 2. Tambah Endpoint CRUD Ruangan di `main.go`

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// Struct Model Data Ruangan
type Room struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Capacity    int    `json:"capacity"`
	IsAvailable bool   `json:"is_available"`
}

func main() {
	// Koneksi DB Sederhana
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=roomify_db sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	// 1. ENDPOINT CREATE (POST /rooms)
	r.POST("/rooms", func(c *gin.Context) {
		var req Room
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
			return
		}

		// Validasi Sederhana
		if req.Name == "" || req.Capacity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan kapasitas wajib diisi valid"})
			return
		}

		req.Id = uuid.New().String()
		req.IsAvailable = true

		// Query Insert ke DB
		query := "INSERT INTO mst_room (id, name, capacity, is_available) VALUES ($1, $2, $3, $4)"
		_, err := db.Exec(query, req.Id, req.Name, req.Capacity, req.IsAvailable)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Ruangan berhasil dibuat!",
			"data":    req,
		})
	})

	// 2. ENDPOINT READ ALL (GET /rooms)
	r.GET("/rooms", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, capacity, is_available FROM mst_room")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var rooms []Room
		for rows.Next() {
			var room Room
			if err := rows.Scan(&room.Id, &room.Name, &room.Capacity, &room.IsAvailable); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			rooms = append(rooms, room)
		}

		c.JSON(http.StatusOK, gin.H{
			"data": rooms,
		})
	})

	r.Run(":8085")
}
```

---

## 3. Uji Coba HTTP Request

1. **POST `/rooms`**
   ```json
   {
     "name": "Ruang Rapat Utama",
     "capacity": 20
   }
   ```
2. **GET `/rooms`** untuk melihat daftar ruangan yang berhasil disimpan.

---

## ⚠️ Masalah Mulai Muncul!
Ketika Anda menambah modul baru (User, Division, Reservation, Facility), file `main.go` Anda akan membengkak hingga ribuan baris!  

Inilah titik di mana **manusia menyadari perlunya Melakukan Refactoring ke Clean Architecture (Fase 4)**.

---

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [Lanjut ke Fase 4 ➡️](04_fase_4_refactoring_clean_architecture.md)
