# 📦 Bagian 3: Model Domain & DTO (`model/`)

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [⬅️ Bagian 2](02_konfigurasi_dan_koneksi_database.md) | [Lanjut ke Bagian 4 ➡️](04_utility_dan_helper.md)

---

## 1. Domain Model (`model/room.go`)

Model mewakili entitas tabel di dalam database.

```go
package model

type Room struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Capacity    int      `json:"capacity"`
	IsAvailable bool     `json:"is_available"`
	RoomTypeId  string   `json:"room_type_id"`
	RoomType    RoomType `json:"room_type"`
}
```

---

## 2. Pagination DTO (`model/dto/paging.go`)

DTO untuk menampung metadata paginasi:

```go
package dto

type Paging struct {
	Page        int `json:"page"`
	RowsPerPage int `json:"rowsPerPage"`
	TotalRows   int `json:"totalRows"`
	TotalPages  int `json:"totalPages"`
}
```

---

## 3. Standard JSON Response (`model/dto/response/json_response.go`)

Format balasan HTTP Response standar untuk menjamin konsistensi payload JSON di seluruh API.

```go
package response

import (
	"go-roomify/model/dto"

	"github.com/gin-gonic/gin"
)

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SingleResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

type PagedResponse struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
	Paging dto.Paging  `json:"paging"`
}

func SendSingleResponse(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(code, SingleResponse{
		Status: Status{Code: code, Message: message},
		Data:   data,
	})
}

func SendPagedResponse(ctx *gin.Context, code int, message string, data interface{}, paging dto.Paging) {
	ctx.JSON(code, PagedResponse{
		Status: Status{Code: code, Message: message},
		Data:   data,
		Paging: paging,
	})
}

func SendSingleResponseError(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, SingleResponse{
		Status: Status{Code: code, Message: message},
		Data:   nil,
	})
}
```

---

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md)
