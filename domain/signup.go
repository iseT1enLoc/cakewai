package domain

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"context"
)

type SignupRequest struct {
	Name     string `form:"name" binding:"required" json:"name"`
	Email    string `form:"email" binding:"required,email" json:"email"`
	Password string `form:"password" binding:"required" json:"password"`
	RoleName string `form:"role_name" bson:"role_name,omitempty" json:"role_name,omitempty"`
	RoleID   string `form:"role_id" bson:"role_id,omitempty" json:"role_id,omitempty"`
}

type SignupResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignupUseCase interface {
	SignUp(ctx context.Context, request SignupRequest, env *appconfig.Env) (accessToken string, refreshToken string, uid string, err error)
}
