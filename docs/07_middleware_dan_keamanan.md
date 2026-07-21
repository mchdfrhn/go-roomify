# 🛡️ Bagian 7: Middleware & Keamanan (`middleware/`)

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md) | [⬅️ Bagian 6](06_business_logic_layer_usecase.md) | [Lanjut ke Bagian 8 ➡️](08_presentation_layer_controller.md)

---

## Auth & Role Guard Middleware (`middleware/auth_middleware.go`)

Middleware bertugas memeriksa ketersediaan token JWT pada *Header Authorization*, mengonfirmasi keabsahan token, serta memvalidasi apakah pengguna memiliki Peran (*Role*) yang diizinkan untuk mengakses suatu endpoint.

```go
package middleware

import (
	"go-roomify/model/dto/response"
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
			response.SendSingleResponseError(ctx, http.StatusUnauthorized, "Unauthorized: Token tidak ditemukan")
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := a.jwtToken.VerifyToken(tokenString)
		if err != nil {
			response.SendSingleResponseError(ctx, http.StatusUnauthorized, "Unauthorized: Token tidak valid")
			ctx.Abort()
			return
		}

		userRole, _ := claims["role"].(string)
		if len(roles) > 0 {
			validRole := false
			for _, r := range roles {
				if r == userRole {
					validRole = true
					break
				}
			}
			if !validRole {
				response.SendSingleResponseError(ctx, http.StatusForbidden, "Forbidden: Akses ditolak")
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

[⬅️ Kembali ke Panduan Utama](../GO_CLEAN_ARCHITECTURE_GUIDE.md)
