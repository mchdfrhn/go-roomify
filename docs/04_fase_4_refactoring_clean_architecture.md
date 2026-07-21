# 🏗️ Fase 4: Refactoring ke Clean Architecture (Memecah Kode)

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [⬅️ Fase 3](03_fase_3_crud_pertama_monolit.md) | [Lanjut ke Fase 5 ➡️](05_fase_5_otentikasi_jwt_dan_middleware.md)

---

## 💡 Mindset Fase 4

File `main.go` di Fase 3 mulai penuh dan berantakan. Saatnya memisahkannya ke dalam **4 Layer Utama**:

1. **`model/`**: Tempat menyimpan bentuk Struct Data.
2. **`repository/`**: Tempat menyimpan kueri SQL/Akses Database.
3. **`usecase/`**: Tempat menyimpan Aturan Bisnis & Validasi.
4. **`delivery/controller/`**: Tempat menangani HTTP Request & Response (Gin Engine).

---

## 1. Pindahkan Struct ke `model/room.go`

```go
package model

type Room struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Capacity    int    `json:"capacity"`
	IsAvailable bool   `json:"is_available"`
}
```

---

## 2. Pindahkan Kueri SQL ke `repository/room_repository.go`

```go
package repository

import (
	"database/sql"
	"go-roomify/model"
)

type RoomRepository interface {
	Create(room model.Room) error
	GetAll() ([]model.Room, error)
}

type roomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(room model.Room) error {
	query := "INSERT INTO mst_room (id, name, capacity, is_available) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, room.Id, room.Name, room.Capacity, room.IsAvailable)
	return err
}

func (r *roomRepository) GetAll() ([]model.Room, error) {
	rows, err := r.db.Query("SELECT id, name, capacity, is_available FROM mst_room")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []model.Room
	for rows.Next() {
		var room model.Room
		if err := rows.Scan(&room.Id, &room.Name, &room.Capacity, &room.IsAvailable); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}
```

---

## 3. Pindahkan Logika Bisnis ke `usecase/room_usecase.go`

```go
package usecase

import (
	"fmt"
	"go-roomify/model"
	"go-roomify/repository"
	"net/http"

	"github.com/google/uuid"
)

type RoomUsecase interface {
	CreateRoom(payload model.Room) (model.Room, int, error)
	FindAllRooms() ([]model.Room, int, error)
}

type roomUsecase struct {
	roomRepo repository.RoomRepository
}

func NewRoomUsecase(roomRepo repository.RoomRepository) RoomUsecase {
	return &roomUsecase{roomRepo: roomRepo}
}

func (u *roomUsecase) CreateRoom(payload model.Room) (model.Room, int, error) {
	if payload.Name == "" || payload.Capacity <= 0 {
		return model.Room{}, http.StatusBadRequest, fmt.Errorf("nama dan kapasitas tidak valid")
	}

	payload.Id = uuid.New().String()
	payload.IsAvailable = true

	if err := u.roomRepo.Create(payload); err != nil {
		return model.Room{}, http.StatusInternalServerError, err
	}

	return payload, http.StatusCreated, nil
}

func (u *roomUsecase) FindAllRooms() ([]model.Room, int, error) {
	rooms, err := u.roomRepo.GetAll()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return rooms, http.StatusOK, nil
}
```

---

## 4. Pindahkan HTTP Handler ke `delivery/controller/room_controller.go`

```go
package controller

import (
	"go-roomify/model"
	"go-roomify/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	uc usecase.RoomUsecase
	rg *gin.RouterGroup
}

func NewRoomController(uc usecase.RoomUsecase, rg *gin.RouterGroup) *RoomController {
	return &RoomController{uc: uc, rg: rg}
}

func (c *RoomController) Route() {
	router := c.rg.Group("/rooms")
	router.POST("", c.createHandler)
	router.GET("", c.getAllHandler)
}

func (c *RoomController) createHandler(ctx *gin.Context) {
	var payload model.Room
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
		return
	}

	data, code, err := c.uc.CreateRoom(payload)
	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(code, gin.H{"message": "Sukses", "data": data})
}

func (c *RoomController) getAllHandler(ctx *gin.Context) {
	data, code, err := c.uc.FindAllRooms()
	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(code, gin.H{"data": data})
}
```

---

## 🎯 Hasil Refactoring

Sekarang kode Anda sudah bersih, terisolasi, mudah di-test, dan tidak ada lagi penggabungan kueri SQL di dalam HTTP Handler!

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [Lanjut ke Fase 5 ➡️](05_fase_5_otentikasi_jwt_dan_middleware.md)
