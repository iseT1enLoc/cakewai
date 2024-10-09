package domain

import (
	"context"

	appconfig "cakewai/cakewai.com/component/appcfg"
)

type RefreshTokenRequest struct {
	RefreshToken string `form:"refreshToken" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenUseCase interface {
	RefreshToken(ctx context.Context, request RefreshTokenRequest, env *appconfig.Env) (accessToken string, refreshToken string, err error)
}
