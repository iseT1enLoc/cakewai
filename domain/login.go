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
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginUseCase interface {
	Login(ctx context.Context, request LoginRequest, env *appconfig.Env) (user *User, accessToken string, refreshToken string, err error)
}
