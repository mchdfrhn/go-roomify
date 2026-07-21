# 🚀 Bagian 1: Inisialisasi Project & Instalasi Dependensi

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md)

---

## 1. Menyiapkan Modul Go

Langkah pertama dalam membuat proyek berbasis Go adalah membuat folder proyek dan menginisialisasi **Go Modules** (`go.mod`). Modul ini berfungsi untuk melacak library pihak ketiga dan dependensi proyek.

```bash
# Buat folder proyek baru
mkdir my-go-project
cd my-go-project

# Inisialisasi modul Go
go mod init my-go-project
```

---

## 2. Instalasi Library Utama

Proyek ini menggunakan arsitektur modular yang didukung oleh beberapa library Go populer:

```bash
# 1. Web Framework (Router & HTTP Engine)
go get -u github.com/gin-gonic/gin

# 2. Driver Database PostgreSQL
go get -u github.com/lib/pq

# 3. Manajer Environment Variable (.env)
go get -u github.com/joho/godotenv

# 4. Autentikasi JSON Web Token (JWT)
go get -u github.com/golang-jwt/jwt/v5

# 5. Generator UUID v4 (Unique Identifier)
go get -u github.com/google/uuid

# 6. Helper Export SQL ke CSV (Opsional)
go get -u github.com/joho/sqltocsv
```

---

## 3. Struktur Direktori Proyek

Setelah inisialisasi, buatlah hirarki folder sesuai dengan **Clean Architecture**:

```bash
mkdir -p config delivery/controller middleware model/dto/request model/dto/response repository usecase utils/common utils/query
```

---

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [Lanjut ke Bagian 2: Konfigurasi & DB ➡️](02_konfigurasi_dan_koneksi_database.md)
