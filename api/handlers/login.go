package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/component/response"
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

		fmt.Println("Entered LoginHandler...")

		// Parse and bind the JSON request
		if err := c.ShouldBindJSON(&request); err != nil {
			fmt.Println("Error parsing login request:", err)
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request payload",
				Error:   err.Error(),
			})
			return
		}
		fmt.Println("Request parsed successfully.")

		// Perform login operation using the use case
		user, accessToken, refreshToken, err := li.LoginUsecase.Login(c, request, li.Env)
		if err != nil {
			fmt.Println("Login failed:", err)
			c.JSON(http.StatusUnauthorized, response.FailedResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid credentials",
				Error:   err.Error(),
			})
			return
		}
		fmt.Println("Login successful. Tokens generated.")
		type loginUserResponse struct {
			AccessToken  string               `json:"access_token"`
			RefreshToken string               `json:"refresh_token"`
			User         *domain.UserResponse `json:"user"`
		}

		// Construct the response
		res := loginUserResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User: &domain.UserResponse{
				Id:             user.Id,
				GoogleId:       user.GoogleId,
				ProfilePicture: user.ProfilePicture,
				Name:           user.Name,
				Email:          user.Email,
				Phone:          user.Phone,
				Address:        user.Address,
				CreatedAt:      user.CreatedAt,
			},
		}

		// Send the response
		c.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Login successful",
			},
			Data: res,
		})
		fmt.Println("LoginHandler completed successfully.")
	}
}
