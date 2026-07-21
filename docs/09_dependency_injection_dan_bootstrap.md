# 🔌 Bagian 9: Dependency Injection & Bootstrap (`delivery/server.go` & `main.go`)

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [⬅️ Bagian 8](08_presentation_layer_controller.md) | [Lanjut ke Bagian 10 ➡️](10_cheat_sheet_fitur_baru.md)

---

## 1. Dependency Injection Container (`delivery/server.go`)

File ini berfungsi menghubungkan seluruh komponen dari layer terbawah hingga terluar (**Database ➔ Repository ➔ Usecase ➔ Controller ➔ Gin Engine**).

```go
package delivery

import (
	"go-roomify/config"
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
	host           string
}

func NewServer() *Server {
	// 1. Load Config & DB Connection
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Config Error: ", err)
	}

	db, err := config.NewDbConnection(cfg)
	if err != nil {
		log.Fatal("DB Connection Error: ", err)
	}

	// 2. Inisialisasi Security & Middleware
	jwtToken := common.NewJwtToken(cfg.TokenConfig)
	authMiddleware := middleware.NewAuthMiddleware(jwtToken)

	// 3. Inisialisasi Repository
	roomRepo := repository.NewRoomRepository(db.Conn())

	// 4. Inisialisasi Usecase
	roomUc := usecase.NewRoomUsecase(roomRepo)

	// 5. Inisialisasi Gin Engine & Router Group
	engine := gin.Default()
	routerGroup := engine.Group("/api/v1")

	return &Server{
		roomUc:         roomUc,
		authMiddleware: authMiddleware,
		engine:         engine,
		routerGroup:    routerGroup,
		host:           cfg.ApiPort,
	}
}

func (s *Server) setupControllers() {
	controller.NewRoomController(s.roomUc, s.routerGroup, s.authMiddleware).Route()
}

func (s *Server) Run() {
	s.setupControllers()
	if err := s.engine.Run(s.host); err != nil {
		log.Fatal("Server Error: ", err)
	}
}
```

---

## 2. Main Entry Point (`main.go`)

Merupakan titik awal eksekusi program. `main()` cukup memanggil `delivery.NewServer().Run()`.

```go
package main

import "go-roomify/delivery"

func main() {
	delivery.NewServer().Run()
}
```

---

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md)
