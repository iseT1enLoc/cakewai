package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/domain"
	"fmt"
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

// "http://localhost:8080/api/public/google/callback",
var googleOauthConfig = &oauth2.Config{
	RedirectURL: "https://www.cakewaibackend.id.vn/api/public/google/callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email",
	},
	Endpoint: google.Endpoint,
}

func (gc *GoogleController) HandleGoogleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(c.Request.URL.RequestURI())
		fmt.Println("Line 35")
		// Generate OAuth state and set it in the cookie
		oauthState := gc.GoogleUseCase.GenerateStateOauthCookie(c.Writer)
		fmt.Println("Line 38")
		// Set Google OAuth configuration
		googleOauthConfig.ClientSecret = gc.Env.GOOGLE_CLIENT_SECRET
		googleOauthConfig.ClientID = gc.Env.GOOGLE_CLIENT_ID
		// Generate the URL to redirect to
		u := googleOauthConfig.AuthCodeURL(oauthState)
		fmt.Println(u)
		fmt.Println(":ine 44")
		// Redirect using Gin's context
		c.Redirect(http.StatusTemporaryRedirect, u)
	}
}
func (gc *GoogleController) HandleGoogleCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(c.Request.URL.RequestURI())

		//redirectURL := os.Getenv("redirect")
		fmt.Println("Enter google callback handler")
		googleOauthConfig.ClientSecret = gc.Env.GOOGLE_CLIENT_SECRET
		googleOauthConfig.ClientID = gc.Env.GOOGLE_CLIENT_ID

		// Get the oauth state cookie
		oauthState, err := c.Cookie("oauthstate")
		fmt.Printf("Enter line 55")
		if err != nil {
			// Handle error if cookie is not found
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
		fmt.Printf("Enter line 61")
		// Validate the state
		if c.Query("state") != oauthState {
			// log.Error("invalid oauth google state")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
		fmt.Printf("Enter line 68")
		// Get user data from Google
		data, err := gc.GoogleUseCase.GetUserDataFromGoogle(googleOauthConfig, c.Query("code"), oauthGoogleUrlAPI)
		if err != nil {
			log.Error(err)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
		fmt.Printf("Enter line 76")
		fmt.Println(data)
		// Perform Google login

		accessToken, refreshToken, err := gc.GoogleUseCase.GoogleLogin(c.Request.Context(), data, gc.Env)
		if err != nil {
			log.Error(err)
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
		// Create the response object
		// loginggresponse := domain.SignupResponse{
		// 	AccessToken:  accessToken,
		// 	RefreshToken: refreshToken,
		// }
		fmt.Println("Enter line 87")
		// Set cookies for access and refresh tokens

		//utils.SetCookie(c.Writer, "access_token", accessToken,false)
		//utils.SetCookie(c.Writer, "refresh_token", refreshToken,true)

		//c.JSON(http.StatusOK, gin.H{"info": "user data"})

		// 	Data:loginggresponse,
		// })
		//redirectURL := "http://localhost:5173/home?token=" + loginggresponse.AccessToken + "&user=" + loginggresponse.User
		// Redirect to the profile page
		fmt.Print("redirect")
		//redirectURL := "http://localhost:5173/home?accesstoken=" + accessToken + "refrestoken=" + refreshToken
		//redirectURL := "http://localhost:5173/home"

		//front_end := os.Getenv("FRONT_END")
		c.Redirect(http.StatusPermanentRedirect, "https://cakewaitown.com/?token="+accessToken+"&refreshToken="+refreshToken)

	}
}
