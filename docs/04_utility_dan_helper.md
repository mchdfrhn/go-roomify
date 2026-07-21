# 🛠️ Bagian 4: Utility & Helper (`utils/`)

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [⬅️ Bagian 3](03_model_domain_dan_dto.md) | [Lanjut ke Bagian 5 ➡️](05_data_access_layer_repository.md)

---

## 1. Paginasi Helper (`utils/paginate.go`)

Fungsi pembantu untuk menghitung halaman (*total pages*) berdasarkan offset dan jumlah total baris.

```go
package utils

import (
	"go-roomify/model/dto"
	"math"
)

func Paginate(page, size, totalRows int) dto.Paging {
	totalPages := int(math.Ceil(float64(totalRows) / float64(size)))
	return dto.Paging{
		Page:        page,
		RowsPerPage: size,
		TotalRows:   totalRows,
		TotalPages:  totalPages,
	}
}
```

---

## 2. JWT Manager (`utils/common/jwt_token.go`)

Helper untuk membuat (*generate*) dan memverifikasi (*verify*) token JWT.

```go
package common

import (
	"fmt"
	"go-roomify/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtToken interface {
	GenerateToken(id, username, role string) (string, error)
	VerifyToken(tokenString string) (jwt.MapClaims, error)
}

type jwtToken struct {
	cfg config.TokenConfig
}

func NewJwtToken(cfg config.TokenConfig) JwtToken {
	return &jwtToken{cfg: cfg}
}

func (j *jwtToken) GenerateToken(id, username, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"role":     role,
		"iss":      j.cfg.JwtIssuer,
		"exp":      time.Now().Add(j.cfg.JwtLifeTime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.cfg.JwtSignature))
}

func (j *jwtToken) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.JwtSignature), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
```

---

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md)
