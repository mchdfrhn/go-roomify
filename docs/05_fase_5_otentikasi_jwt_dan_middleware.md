# 🔑 Fase 5: Otentikasi JWT & Middleware Protection

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [⬅️ Fase 4](04_fase_4_refactoring_clean_architecture.md) | [Lanjut ke Fase 6 ➡️](06_fase_6_dependency_injection_dan_scaleup.md)

---

## 💡 Mindset Fase 5

Setelah fitur utama terpisah dengan rapi, aplikasi Anda siap diberi lapisan **Keamanan & Autentikasi** untuk membatasi hak akses pengguna (RBAC: Admin vs User).

---

## 1. Install Package JWT

```bash
go get -u github.com/golang-jwt/jwt/v5
```

---

## 2. Buat JWT Helper (`utils/common/jwt_token.go`)

```go
package common

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtToken interface {
	GenerateToken(id, username, role string) (string, error)
	VerifyToken(tokenString string) (jwt.MapClaims, error)
}

type jwtToken struct {
	secretKey []byte
}

func NewJwtToken(secretKey string) JwtToken {
	return &jwtToken{secretKey: []byte(secretKey)}
}

func (j *jwtToken) GenerateToken(id, username, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *jwtToken) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token tidak valid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("claims tidak valid")
	}
	return claims, nil
}
```

---

## 3. Buat Middleware Role Guard (`middleware/auth_middleware.go`)

```go
package middleware

import (
	"go-roomify/utils/common"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	RequireToken(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtToken common.JwtToken
}

func NewAuthMiddleware(jwtToken common.JwtToken) AuthMiddleware {
	return &authMiddleware{jwtToken: jwtToken}
}

func (a *authMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := a.jwtToken.VerifyToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired / invalid"})
			ctx.Abort()
			return
		}

		userRole, _ := claims["role"].(string)
		if len(roles) > 0 {
			valid := false
			for _, r := range roles {
				if r == userRole {
					valid = true
					break
				}
			}
			if !valid {
				ctx.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak untuk peran Anda"})
				ctx.Abort()
				return
			}
		}

		ctx.Set("user_id", claims["id"])
		ctx.Next()
	}
}
```

---

## 4. Pasang Middleware di Controller (`delivery/controller/room_controller.go`)

```go
func (c *RoomController) Route() {
	router := c.rg.Group("/rooms")

	// Hanya ADMIN yang boleh membuat ruangan baru
	router.POST("", c.authMiddleware.RequireToken("ADMIN"), c.createHandler)

	// ADMIN dan USER boleh melihat daftar ruangan
	router.GET("", c.authMiddleware.RequireToken("ADMIN", "USER"), c.getAllHandler)
}
```

---

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [Lanjut ke Fase 6 ➡️](06_fase_6_dependency_injection_dan_scaleup.md)
