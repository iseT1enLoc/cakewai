package domain

import (
	"github.com/golang-jwt/jwt/v4"
)

type JwtCustomClaims struct {
	Name     string `json:"name"`
	ID       string `json:"_id"`
	Email    string `json:"email"`
	GoogleId string `json:"google_id"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	Name     string `json:"name"`
	ID       string `json:"_id"`
	Email    string `json:"email"`
	GoogleId string `json:"google_id"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}
