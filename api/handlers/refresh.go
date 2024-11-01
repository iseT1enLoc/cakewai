package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/component/response"
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
		fmt.Println(request)
		accessToken, _, err := u.RefreshTokenUsecase.RefreshToken(c, request, request.RefreshToken, u.Env)
		fmt.Println("line 28 refresh token handler")
		fmt.Println(accessToken)
		fmt.Println(request.RefreshToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response := domain.RefreshShortResponse{
			AccessToken:  accessToken,
			RefreshToken: request.RefreshToken,
		}

		c.JSON(http.StatusOK, response)
	}
}
func (u *RefreshTokenHandler) RevokeRefreshTokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request domain.RefreshTokenRequest

		//Bind JSON request body to RefreshTokenRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			fmt.Println("Error at line 23")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(request)

		err := u.RefreshTokenUsecase.RevokeToken(c, request.RefreshToken, u.Env)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response.BasicResponse{
			Code:    200,
			Message: "You have successfully revoke refresh token",
			Error:   "None",
			Data:    "okokkokok",
		})
	}
}
