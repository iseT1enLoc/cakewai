package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/domain"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ydb-platform/ydb-go-sdk/v3/log"
)

type SignupController struct {
	SignupUseCase domain.SignupUseCase
	Env           *appconfig.Env
}

func (sc *SignupController) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request domain.SignupRequest

		fmt.Print("line 22 at sign up handler")
		// Decode the JSON request body
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Print("line 30 at signup handler")
		// Call the signup use case
		accessToken, refreshToken, err := sc.SignupUseCase.SignUp(c.Request.Context(), request, sc.Env)
		if err != nil {
			log.Error(err)
			fmt.Println("Error happened in line 35 sign up handler")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already existed"})
			return
		}

		fmt.Print("line 39 at sign up handler")
		// Create the response object
		signupResponse := domain.SignupResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		fmt.Print("line 45 sign up handler")
		// Send the response
		c.JSON(http.StatusOK, signupResponse)
	}
}
