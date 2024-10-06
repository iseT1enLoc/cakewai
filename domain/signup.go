package domain

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"context"
)

type SignupRequest struct {
	Name     string `form:"name" binding:"required"`
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type SignupResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignupUseCase interface {
	SignUp(ctx context.Context, request SignupRequest, env *appconfig.Env) (accessToken string, refreshToken string, err error)
}
