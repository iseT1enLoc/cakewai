package domain

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JwtCustomClaims struct {
	Name     string    `json:"name"`
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	GoogleId string    `json:"google_id"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	Name     string    `json:"name"`
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	GoogleId string    `json:"google_id"`
	jwt.RegisteredClaims
}
