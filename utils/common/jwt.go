package common

import (
	"time"
	"errors"
	"go-roomify/config"
	"go-roomify/model"

	"github.com/dgrijalva/jwt-go"
)

type JwtClaim struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
	Role string `json:"role"`
	//Email string `json:"email"`
}

type JwtToken interface {
	GenerateTokenJwt(user_data model.UserCredentialJwt) (string, error)
	VerifyToken(token_string string) (jwt.MapClaims, error)
}

type jwtToken struct {
	config config.TokenConfig
}

func (self *jwtToken) GenerateTokenJwt(user_cr model.UserCredentialJwt) (string, error) {
	claims := JwtClaim{
		StandardClaims: jwt.StandardClaims{
			Issuer:    self.config.IssuerName,
			ExpiresAt: time.Now().Add(self.config.JwtLifeTime).Unix(),
		},
		UserId: user_cr.Id,
		Role: user_cr.Role,
		//Email: user_data.Email
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed_token, err := token.SignedString(self.config.JwtSignatureKey)

	if err != nil {
		return "", err
	}
	return signed_token, nil
}

func (self *jwtToken) VerifyToken(token_string string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(token_string, func(token *jwt.Token) (any, error) {
		return self.config.JwtSignatureKey, nil
	})
	if err != nil {
		return nil, err
	}

	// COnvert jwt claims ke jwt map claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Failed to parse map claims or token is not valid")
	}

	// Verifikasi issuer
	if !claims.VerifyIssuer(self.config.IssuerName, true) {
		return nil, errors.New("Failed to Verify Issuer Name")
	}

	// Verifikasi expired
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, errors.New("Token is expired")
	}

	return claims, nil
}

func NewJwtToken(token_config config.TokenConfig) JwtToken {
	return &jwtToken{
		config: token_config,
	}
}