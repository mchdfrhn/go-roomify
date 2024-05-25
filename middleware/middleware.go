package middleware

import (
	// "go-roomify/config"
	// "go-roomify/model"
	"go-roomify/utils/common"
	"go-roomify/model/dto/response"
	"net/http"
	//"os"
	"strings"
	//"time"
	//"fmt"

	//"github.com/golang-jwt/jwt/v5"
	//"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	RequireToken(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtToken common.JwtToken
}

type authHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}

func (self *authMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var auth_header authHeader

		if err := ctx.ShouldBindHeader(&auth_header); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		// auth_header := ctx.GetHeader("Authorization")

		//fmt.Println("HEADER AUTH: ", auth_header)

		// if !strings.Contains(auth_header, "Bearer") {
		// 	response.SendSingleResponseError(
		// 		ctx,
		// 		http.StatusBadRequest,
		// 		"Invalid Token",
		// 	)
		// 	return
		// }
		// jwt_signature_key := []byte(os.Getenv("SIGNATURE"))
		token_string := strings.Replace(auth_header.AuthorizationHeader, "Bearer ", "", -1)
		if token_string == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Status{
				Code: http.StatusUnauthorized,
				Description: "Unauthorized",
			})
			return
		}
		// jwt_claims := &common.JwtClaim{}

		// token, err := jwt.ParseWithClaims(token_string, jwt_claims, func(token *jwt.Token) (any, error) {
		// 	return jwt_signature_key, nil
		// })

		jwt_claims, err := self.jwtToken.VerifyToken(token_string)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Status{
				Code: http.StatusUnauthorized,
				Description: err.Error(),
			})
			return
		}

		// if !token.Valid {
		// 	response.SendSingleResponseError(
		// 		ctx,
		// 		http.StatusUnauthorized,
		// 		"Unauthorized User",
		// 	)
		// 	return
		// }

		// expired_at := jwt_claims.ExpiresAt

		// fmt.Println("EXPIRED AT: ", expired_at)
		// fmt.Println("CLAIMS: ", jwt_claims)

		// if time.Now().Unix() > expired_at {
		// 	response.SendSingleResponseError(
		// 		ctx,
		// 		http.StatusUnauthorized,
		// 		"Expired Token",
		// 	)
		// 	return
		// }

		valid_role := false

		if len(roles) > 0 {
			for _, role := range roles {
				if role == jwt_claims["role"] { //jwt_claims.UserData.Role {
					valid_role = true
					break
				}
			}
		}

		if !valid_role {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Status{
				Code: http.StatusUnauthorized,
				Description: "You don't have permission",
			})
			return
		}

		ctx.Next()
	}
}

func NewAuthMiddleware(jwtToken common.JwtToken) AuthMiddleware {
	return &authMiddleware {
		jwtToken: jwtToken,
	}
}