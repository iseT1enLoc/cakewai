package domain

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountRecovery interface {
	ResetPasswordRequest(ctx context.Context, env *appconfig.Env, email string) error
	ResetPasswordProcessing(ctx context.Context, env *appconfig.Env, new_password string, token string) error
	ChangesPassword(ctx context.Context, env *appconfig.Env, userId primitive.ObjectID, currentPassword string, newPassword string) error
}

type ResetPasswordReq struct {
	Email string `json:"email" bson:"email"`
}
