# 🌱 Fase 1: Hello World & HTTP Server Paling Sederhana

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [Lanjut ke Fase 2 ➡️](02_fase_2_koneksi_database.md)

---

## 💡 Mindset Fase 1

Jangan langsung membuat puluhan folder! Di awal pembuatan aplikasi, fokus kita hanya 1:  
**Memastikan proyek Go dapat di-*run* dan HTTP Server bisa menerima request dari luar.**

---

## 1. Inisialisasi Project

Buka terminal di komputer Anda:

```bash
mkdir my-roomify
cd my-roomify

# Buat modul Go
go mod init my-roomify

# Download framework Gin Gonic (Web Framework ringan)
go get -u github.com/gin-gonic/gin
```

---

## 2. Tulis `main.go` Pertama Anda (Cuma 15 Baris!)

Buat file `main.go` di folder root:

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi router bawaan Gin
	r := gin.Default()

	// Tes endpoint paling dasar (Ping Check)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"status":  "server up and running",
		})
	})

	// Jalankan server di port 8085
	r.Run(":8085")
}
```

---

## 3. Jalankan dan Tes Server

Jalankan perintah ini di terminal:

```bash
go run main.go
```

Buka browser atau Postman, lalu akses:
```text
http://localhost:8085/ping
```

**Hasil Response JSON:**
```json
{
  "message": "pong",
  "status": "server up and running"
}
```

✅ **Selamat!** Server web pertama Anda sudah berjalan dengan sukses. Sekarang saatnya melangkah ke **Fase 2: Menghubungkan ke Database**.

---

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [Lanjut ke Fase 2 ➡️](02_fase_2_koneksi_database.md)
