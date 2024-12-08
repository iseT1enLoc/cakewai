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

// func (li *LoginHandler) LoginHandler() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var request domain.LoginRequest
// 		// Context timeout handling
// 		ctx, cancel := context.WithTimeout(c, time.Second*5)
// 		defer cancel()

// 		fmt.Printf("line 21 login handler")
// 		// Bind JSON request body to LoginRequest
// 		if err := c.ShouldBindJSON(&request); err != nil {
// 			c.JSON(http.StatusBadRequest, response.FailedResponse{
// 				Code:    http.StatusBadRequest,
// 				Message: "Can not parsing the resquest",
// 				Error:   err.Error(),
// 			})
// 			return
// 		}
// 		fmt.Print("line 27 login handler")

// 		accessToken, refreshToken, err := li.LoginUsecase.Login(c, request, li.Env)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, response.FailedResponse{
// 				Code:    http.StatusBadRequest,
// 				Message: "Can not parsing the resquest",
// 				Error:   err.Error(),
// 			})
// 			return
// 		}
// 		fmt.Print("line 34 login handler")
// 		loginresponse := domain.LoginResponse{
// 			AccessToken:  accessToken,
// 			RefreshToken: refreshToken,
// 		}
// 		// Set user ID in context
// 		// Send the response
// 		c.JSON(http.StatusOK, response.Success{
// 			ResponseFormat: response.ResponseFormat{
// 				Code:    http.StatusCreated,
// 				Message: "Login Successfully",
// 			},
// 			Data: loginresponse,
// 		})

//		}
//	}
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
		accessToken, refreshToken, err := li.LoginUsecase.Login(c, request, li.Env)
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

		// Construct the response structure
		loginResponse := domain.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		// Send the response
		c.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Login successful",
			},
			Data: loginResponse,
		})
		fmt.Println("LoginHandler completed successfully.")
	}
}
