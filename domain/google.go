package domain

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type GoogleUser struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type GoogleUseCase interface {
	GoogleLogin(ctx context.Context, data []byte, env *appconfig.Env) (accessToken string, refreshToken string, err error)
	GetUserDataFromGoogle(googleOauthConfig *oauth2.Config, code, oauthGoogleUrlAPI string) ([]byte, error)
	GenerateStateOauthCookie(w http.ResponseWriter) string
}
