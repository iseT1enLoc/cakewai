package domain

import (
	"context"
	"time"

	appconfig "cakewai/cakewai.com/component/appcfg"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshTokenRequest struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	RefreshToken string             `form:"refresh_token" binding:"required" json:"refresh_token" bson:"refresh_token"`
	UserID       string             `json:"user_id" bson:"user_id"`
	ExpireAt     time.Time          `json:"expire_at" bson:"expire_at"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	RevokeAt     time.Time          `json:"revoke_at" bson:"revoke_at"`
	ReplaceByRT  string             `json:"replaced_token" bson:"replace_token"`
	IsActive     bool               `json:"is_active" bson:"is_active"`
	IsExpire     bool               `json:"is_expire" bson:"is_expire"`
}

type RefreshTokenResponse struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	AccessToken  string             `json:"accessToken"`
	RefreshToken string             `json:"refreshToken"`
	UserID       string             `json:"user_id" bson:"user_id"`
	ExpireAt     time.Time          `json:"expire_at" bson:"expire_at"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	RevokeAt     time.Time          `json:"revoke_at" bson:"revoke_at"`
	ReplaceByRT  string             `json:"replaced_token" bson:"replace_token"`
	IsActive     bool               `json:"is_active" bson:"is_active"`
	IsExpire     bool               `json:"is_expire" bson:"is_expire"`
}

type RefreshTokenUseCase interface {
	RefreshToken(ctx context.Context, request RefreshTokenRequest, env *appconfig.Env) (accessToken string, refreshToken string, err error)
}
