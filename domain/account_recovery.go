package domain

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"context"
)

type AccountRecovery interface {
	ResetPasswordRequest(ctx context.Context, env *appconfig.Env, email string) error
	ResetPasswordProcessing(ctx context.Context, env *appconfig.Env, new_password string, token string) error
}

type ResetPasswordReq struct {
	Email string `json:"email" bson:"email"`
}
