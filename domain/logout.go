package domain

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"context"
)

type LogoutRequest struct {
	Refresh_token string `json:"refresh_token"`
}

type LogOutUseCase interface {
	LogOut(ctx context.Context, request LogoutRequest, env *appconfig.Env) (err error)
}
