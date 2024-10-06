package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/domain"
	"cakewai/cakewai.com/internals/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ydb-platform/ydb-go-sdk/v3/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type GoogleController struct {
	GoogleUseCase domain.GoogleUseCase
	Env           *appconfig.Env
}

var googleOauthConfig = &oauth2.Config{
	RedirectURL: "http://localhost:8080/api/google/callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
	},
	Endpoint: google.Endpoint,
}

func (gc *GoogleController) HandleGoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate OAuth state and set it in the cookie
		oauthState := gc.GoogleUseCase.GenerateStateOauthCookie(c.Writer)

		// Set Google OAuth configuration
		googleOauthConfig.ClientSecret = gc.Env.GOOGLE_CLIENT_SECRET
		googleOauthConfig.ClientID = gc.Env.GOOGLE_CLIENT_ID
		// Generate the URL to redirect to
		u := googleOauthConfig.AuthCodeURL(oauthState)

		// Redirect using Gin's context
		c.Redirect(http.StatusTemporaryRedirect, u)
	}
}
func (gc *GoogleController) HandleGoogleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		googleOauthConfig.ClientSecret = gc.Env.GOOGLE_CLIENT_SECRET
		googleOauthConfig.ClientID = gc.Env.GOOGLE_CLIENT_ID

		// Get the oauth state cookie
		oauthState, err := c.Cookie("oauthstate")
		if err != nil {
			// Handle error if cookie is not found
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		// Validate the state
		if c.Query("state") != oauthState {
			// log.Error("invalid oauth google state")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		// Get user data from Google
		data, err := gc.GoogleUseCase.GetUserDataFromGoogle(googleOauthConfig, c.Query("code"), oauthGoogleUrlAPI)
		if err != nil {
			log.Error(err)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		// Perform Google login
		accessToken, refreshToken, err := gc.GoogleUseCase.GoogleLogin(c.Request.Context(), data, gc.Env)
		if err != nil {
			log.Error(err)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		// Set cookies for access and refresh tokens
		utils.SetCookie(c.Writer, "access_token", accessToken)
		utils.SetCookie(c.Writer, "refresh_token", refreshToken)

		// Redirect to the profile page
		c.Redirect(http.StatusTemporaryRedirect, "/profile")
	}
}
