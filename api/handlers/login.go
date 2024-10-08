package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/domain"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	LoginUsecase domain.LoginUseCase
	Env          *appconfig.Env
}

func (li *LoginHandler) LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request domain.LoginRequest

		fmt.Printf("line 21 login handler")
		// Bind JSON request body to LoginRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"erorr": err})
			return
		}
		fmt.Print("line 27 login handler")

		accessToken, refreshToken, err := li.LoginUsecase.Login(c, request, li.Env)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		fmt.Print("line 34 login handler")
		response := domain.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		fmt.Print("line 39 login handler")
		c.JSON(http.StatusOK, response)
	}
}
