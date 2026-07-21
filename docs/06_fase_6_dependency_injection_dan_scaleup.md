# 🚀 Fase 6: Dependency Injection & Scale-Up Sistem

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [⬅️ Fase 5](05_fase_5_otentikasi_jwt_dan_middleware.md)

---

## 💡 Mindset Fase 6

Pada tahap akhir ini, semua komponen dihubungkan melalui **Dependency Injection Container** (`delivery/server.go`).  
Hal ini membuat `main.go` kembali sangat ringkas (hanya 2-3 baris) dan aplikasi siap untuk di-*scale up* ke puluhan modul baru!

---

## 1. Pembuatan Dependency Injection Container (`delivery/server.go`)

```go
package delivery

import (
	"database/sql"
	"go-roomify/delivery/controller"
	"go-roomify/middleware"
	"go-roomify/repository"
	"go-roomify/usecase"
	"go-roomify/utils/common"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	roomUc         usecase.RoomUsecase
	authMiddleware middleware.AuthMiddleware
	engine         *gin.Engine
	routerGroup    *gin.RouterGroup
}

func NewServer() *Server {
	// 1. Koneksi DB
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=roomify_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// 2. Init Security & Utils
	jwtToken := common.NewJwtToken("SecretKeySangatRahasia123!")
	authMiddleware := middleware.NewAuthMiddleware(jwtToken)

	// 3. Init Repositories
	roomRepo := repository.NewRoomRepository(db)

	// 4. Init Usecases
	roomUc := usecase.NewRoomUsecase(roomRepo)

	// 5. Init Router Engine
	engine := gin.Default()
	routerGroup := engine.Group("/api/v1")

	return &Server{
		roomUc:         roomUc,
		authMiddleware: authMiddleware,
		engine:         engine,
		routerGroup:    routerGroup,
	}
}

func (s *Server) setupControllers() {
	controller.NewRoomController(s.roomUc, s.routerGroup, s.authMiddleware).Route()
}

func (s *Server) Run() {
	s.setupControllers()
	s.engine.Run(":8085")
}
```

---

## 2. Main Entry Point Bersih (`main.go`)

```go
package main

import "go-roomify/delivery"

func main() {
	// Panggil server container dan jalankan
	delivery.NewServer().Run()
}
```

---

## 📈 Trik Menambah Modul Baru (Contoh: Modul `Facility`)

Bila ingin menambah modul baru, Anda tinggal mengulang 4 file sederhana berikut:
1. `model/facility.go`
2. `repository/facility_repository.go`
3. `usecase/facility_usecase.go`
4. `delivery/controller/facility_controller.go`
5. Daftarkan di `delivery/server.go`.

---

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md)
