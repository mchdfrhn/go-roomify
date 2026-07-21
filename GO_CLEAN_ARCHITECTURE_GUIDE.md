# 📘 Panduan Utama Arsitektur Project (`go-roomify`)

Selamat datang di Panduan Arsitektur & Blueprint pembuatan aplikasi **`go-roomify`**. Dokumentasi ini disusun secara terstruktur ke dalam folder `docs/` untuk memudahkan Anda mempelajari dan menggunakannya sebagai acuan pembuatan proyek Go berbasis **Clean Architecture**.

---

## 🗂️ Indeks Dokumentasi Detail (`docs/`)

Setiap tahapan dan modul pembuatan proyek telah dipecah ke dalam file dokumentasi terpisah di bawah ini:

| Bagian | Topik & Modul | Deskripsi | Link Dokumentasi |
| :--- | :--- | :--- | :--- |
| **01** | **Inisialisasi & Dependensi** | Langkah awal `go mod init` dan instalasi library utama (`gin`, `pq`, `godotenv`, `jwt`, `uuid`). | [📄 Buka Panduan](docs/01_inisialisasi_dan_dependensi.md) |
| **02** | **Konfigurasi & Database** | Pengaturan file `.env`, pembacaan file environment (`config/config.go`), dan koneksi database (`config/db_connection.go`). | [📄 Buka Panduan](docs/02_konfigurasi_dan_koneksi_database.md) |
| **03** | **Model & DTO** | Penyiapan entitas domain model, DTO Paginasi, dan skema balasan HTTP response terstandarisasi (`model/dto/response`). | [📄 Buka Panduan](docs/03_model_domain_dan_dto.md) |
| **04** | **Utility & Helper** | Fungsi pembantu umum seperti kalkulator paginasi (`utils/paginate.go`) dan JWT Token Generator/Verifier (`utils/common`). | [📄 Buka Panduan](docs/04_utility_dan_helper.md) |
| **05** | **Data Access Layer** | Pembuatan repository interface & eksekusi query SQL dengan `database/sql` & Custom Query Builder (`repository/`). | [📄 Buka Panduan](docs/05_data_access_layer_repository.md) |
| **06** | **Business Logic Layer** | Implementasi aturan bisnis, validasi, orchestration data, dan mapping HTTP status code (`usecase/`). | [📄 Buka Panduan](docs/06_business_logic_layer_usecase.md) |
| **07** | **Middleware & Keamanan** | Proteksi endpoint menggunakan JWT token verification dan Role-Based Access Control (RBAC) di `middleware/`. | [📄 Buka Panduan](docs/07_middleware_dan_keamanan.md) |
| **08** | **Presentation Layer** | Penyiapan Gin HTTP Controllers, penanganan payload request, routing, dan pemanggilan usecase (`delivery/controller/`). | [📄 Buka Panduan](docs/08_presentation_layer_controller.md) |
| **09** | **Dependency Injection** | Perangkaian hubungan antar layer (**DB ➔ Repo ➔ Usecase ➔ Controller ➔ Router**) di `delivery/server.go` dan `main.go`. | [📄 Buka Panduan](docs/09_dependency_injection_dan_bootstrap.md) |
| **10** | **Cheat Sheet Fitur Baru** | Panduan cepat langkah-demi-langkah (6 urutan baku) ketika hendak menambah modul/fitur baru ke dalam proyek. | [📄 Buka Panduan](docs/10_cheat_sheet_fitur_baru.md) |

---

## 🏗️ Ringkasan Struktur Folder Proyek

```text
go-roomify/
├── GO_CLEAN_ARCHITECTURE_GUIDE.md   # Hub Utama Panduan (File Ini)
├── docs/                            # Folder Dokumentasi Modul Terpisah
│   ├── 01_inisialisasi_dan_dependensi.md
│   ├── 02_konfigurasi_dan_koneksi_database.md
│   ├── 03_model_domain_dan_dto.md
│   ├── 04_utility_dan_helper.md
│   ├── 05_data_access_layer_repository.md
│   ├── 06_business_logic_layer_usecase.md
│   ├── 07_middleware_dan_keamanan.md
│   ├── 08_presentation_layer_controller.md
│   ├── 09_dependency_injection_dan_bootstrap.md
│   └── 10_cheat_sheet_fitur_baru.md
├── config/                          # Konfigurasi & DB Connection
├── delivery/                        # Presentation Layer (Controllers & Server)
├── middleware/                      # HTTP Middleware (Auth JWT)
├── model/                           # Domain Entities & DTOs
├── repository/                      # Data Access Layer (SQL Queries)
├── usecase/                         # Business Logic Layer
├── utils/                           # Helper Functions & Query Builder
├── .env                             # Environment Variables
├── go.mod                           # Go Module Spec
└── main.go                          # Application Entry Point
```
