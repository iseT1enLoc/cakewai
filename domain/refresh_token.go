package domain

import (
	"context"
	"time"

	appconfig "cakewai/cakewai.com/component/appcfg"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshTokenRequest struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	RefreshToken string             ` json:"refresh_token" bson:"refresh_token" form:"refresh_token" binding:"required"`
	UserID       string             `json:"user_id" bson:"user_id"`
	ExpireAt     time.Time          `json:"expire_at" bson:"expire_at"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	RevokeAt     time.Time          `json:"revoke_at" bson:"revoke_at"`
	ReplaceByRT  string             `json:"replaced_token" bson:"replace_token"`
	IsActive     bool               `json:"is_active" bson:"is_active"`
	IsExpire     bool               `json:"is_expire" bson:"is_expire"`
}

type RefreshTokenResponse struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	AccessToken  string             `json:"access_token" bson:"access_token"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	UserID       string             `json:"user_id" bson:"user_id,omitempty"`
	ExpireAt     time.Time          `json:"expire_at" bson:"expire_at,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at,omitempty"`
	RevokeAt     time.Time          `json:"revoke_at" bson:"revoke_at,omitempty"`
	ReplaceByRT  string             `json:"replaced_token" bson:"replace_token,omitempty"`
	IsActive     bool               `json:"is_active" bson:"is_active,omitempty"`
	IsExpire     bool               `json:"is_expire" bson:"is_expire,omitempty"`
}
type RefreshShortResponse struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}
type RefreshTokenUseCase interface {
	//renew access token
	RefreshToken(ctx context.Context, request RefreshTokenRequest, is_admin bool, env *appconfig.Env) (accessToken string, refreshToken string, err error)
	//thu hoi refresh token
	RevokeToken(ctx context.Context, current_RT string, env *appconfig.Env) error
	InsertRefreshTokenToDB(ctx context.Context, refresh_token RefreshTokenRequest, user_id string, is_admin bool, env *appconfig.Env) (string, error)
	GetRefreshTokenFromDB(ctx context.Context, current_refresh_token string, env *appconfig.Env) (*RefreshTokenRequest, error)
	UpdateRefreshTokenChanges(ctx context.Context, updatedRT RefreshTokenRequest, env *appconfig.Env) (*RefreshTokenRequest, error)
	RenewAccessToken(ctx context.Context, refresh RefreshTokenRequest, env *appconfig.Env) (access_token string, refresh_token string, err error)
	DeleteRefreshtoken(ctx context.Context, current_RT string, env *appconfig.Env) error
}
