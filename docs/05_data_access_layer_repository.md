# 🗄️ Bagian 5: Data Access Layer (`repository/`)

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [⬅️ Bagian 4](04_utility_dan_helper.md) | [Lanjut ke Bagian 6 ➡️](06_business_logic_layer_usecase.md)

---

## Implementasi Repository (`repository/room_repository.go`)

Layer Repository mengisolasi kueri SQL dan interaksi database sehingga logika bisnis tidak bergantung pada detail penyimpanan data.

```go
package repository

import (
	"database/sql"
	"go-roomify/model"
	"go-roomify/utils/query"
)

// Interface kontrak
type RoomRepository interface {
	Create(room model.Room) error
	GetById(id string) (model.Room, error)
}

// Implementasi konkrit
type roomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(room model.Room) error {
	q := query.QInsert{DB: r.db}
	_, err := q.Table("mst_room").
		Column("id", "name", "capacity", "is_available", "room_type_id").
		Values(room.Id, room.Name, room.Capacity, room.IsAvailable, room.RoomTypeId).
		Run()
	return err
}

func (r *roomRepository) GetById(id string) (model.Room, error) {
	q := query.QSelect{DB: r.db}
	var room model.Room

	err := q.Table("mst_room").
		Column("id", "name", "capacity", "is_available", "room_type_id").
		Where("id", "=", id).
		RunRow().Scan(&room.Id, &room.Name, &room.Capacity, &room.IsAvailable, &room.RoomTypeId)

	if err != nil {
		return model.Room{}, err
	}
	return room, nil
}
```

---

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md)
