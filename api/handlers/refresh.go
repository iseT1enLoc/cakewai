package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/domain"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RefreshTokenHandler struct {
	RefreshTokenUsecase domain.RefreshTokenUseCase
	Env                 *appconfig.Env
}

func (u *RefreshTokenHandler) RefreshTokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request domain.RefreshTokenRequest

		//Bind JSON request body to RefreshTokenRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			fmt.Println("Error at line 23")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		accessToken, refreshToken, err := u.RefreshTokenUsecase.RefreshToken(c, request, u.Env)
		fmt.Println("line 28 refresh token handler")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response := domain.RefreshTokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		c.JSON(http.StatusOK, response)
	}
}
