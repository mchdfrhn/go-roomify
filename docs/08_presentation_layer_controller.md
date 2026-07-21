# 🌐 Bagian 8: Presentation Layer (`delivery/controller/`)

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [⬅️ Bagian 7](07_middleware_dan_keamanan.md) | [Lanjut ke Bagian 9 ➡️](09_dependency_injection_dan_bootstrap.md)

---

## HTTP Handler (`delivery/controller/room_controller.go`)

Controller bertindak sebagai jembatan antara protokol HTTP dan logika aplikasi (Usecase). Controller membaca query param/body JSON, memanggil usecase, dan mengembalikan format balasan JSON menggunakan DTO Response.

```go
package controller

import (
	"go-roomify/middleware"
	"go-roomify/model"
	"go-roomify/model/dto/response"
	"go-roomify/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	roomUc         usecase.RoomUsecase
	rg             *gin.RouterGroup
	authMiddleware middleware.AuthMiddleware
}

func NewRoomController(roomUc usecase.RoomUsecase, rg *gin.RouterGroup, authMiddleware middleware.AuthMiddleware) *RoomController {
	return &RoomController{
		roomUc:         roomUc,
		rg:             rg,
		authMiddleware: authMiddleware,
	}
}

func (c *RoomController) Route() {
	router := c.rg.Group("/rooms")

	// Pendaftaran Endpoint & Middleware Proteksi Role
	router.POST("", c.authMiddleware.RequireToken("ADMIN"), c.createHandler)
	router.GET("/:id", c.authMiddleware.RequireToken("ADMIN", "USER"), c.getByIdHandler)
}

func (c *RoomController) createHandler(ctx *gin.Context) {
	var payload model.Room
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		response.SendSingleResponseError(ctx, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	data, code, err := c.roomUc.CreateNewRoom(payload)
	if err != nil {
		response.SendSingleResponseError(ctx, code, err.Error())
		return
	}

	response.SendSingleResponse(ctx, code, "Ruangan berhasil ditambahkan", data)
}

func (c *RoomController) getByIdHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	data, code, err := c.roomUc.FindRoomById(id)
	if err != nil {
		response.SendSingleResponseError(ctx, code, err.Error())
		return
	}

	response.SendSingleResponse(ctx, code, "Success", data)
}
```

---

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md)
