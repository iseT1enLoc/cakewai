package domain

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"context"
)

type LoginRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginUseCase interface {
	Login(ctx context.Context, request LoginRequest, env *appconfig.Env) (accessToken string, refreshToken string, err error)
}
