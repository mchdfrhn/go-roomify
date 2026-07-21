# 🧠 Bagian 6: Business Logic Layer (`usecase/`)

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [⬅️ Bagian 5](05_data_access_layer_repository.md) | [Lanjut ke Bagian 7 ➡️](07_middleware_dan_keamanan.md)

---

## Implementasi Usecase (`usecase/room_usecase.go`)

Layer Usecase bertanggung jawab atas seluruh Aturan Bisnis (*Business Rules*), pemrosesan data, validasi, dan koordinasi antar repository. Usecase juga mengembalikan HTTP Status Code yang sesuai untuk diproses oleh controller.

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
	CreateNewRoom(payload model.Room) (model.Room, int, error)
	FindRoomById(id string) (model.Room, int, error)
}

type roomUsecase struct {
	roomRepo repository.RoomRepository
}

func NewRoomUsecase(roomRepo repository.RoomRepository) RoomUsecase {
	return &roomUsecase{roomRepo: roomRepo}
}

func (u *roomUsecase) CreateNewRoom(payload model.Room) (model.Room, int, error) {
	// 1. Validasi Input Bisnis
	if payload.Name == "" || payload.Capacity <= 0 {
		return model.Room{}, http.StatusBadRequest, fmt.Errorf("nama dan kapasitas ruangan wajib diisi valid")
	}

	// 2. Generate UUID & State Awal
	payload.Id = uuid.New().String()
	payload.IsAvailable = true

	// 3. Simpan ke Database via Repository
	err := u.roomRepo.Create(payload)
	if err != nil {
		return model.Room{}, http.StatusInternalServerError, err
	}

	return payload, http.StatusCreated, nil
}

func (u *roomUsecase) FindRoomById(id string) (model.Room, int, error) {
	room, err := u.roomRepo.GetById(id)
	if err != nil {
		return model.Room{}, http.StatusNotFound, fmt.Errorf("ruangan tidak ditemukan")
	}
	return room, http.StatusOK, nil
}
```

---

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md)
