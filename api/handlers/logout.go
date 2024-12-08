package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/domain"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogoutHandler struct {
	LogoutUsecase domain.LogOutUseCase
	Env           *appconfig.Env
}

func (li *LogoutHandler) LogoutHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var logoutRequest domain.LogoutRequest

		// Bind and validate the request
		if err := ctx.ShouldBindJSON(&logoutRequest); err != nil {
			fmt.Printf("Invalid logout request: %v", err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request body",
				Error:   err.Error(),
			})
			return
		}

		// Call the Logout use case
		err := li.LogoutUsecase.LogOut(ctx.Request.Context(), logoutRequest, li.Env)
		if err != nil {
			fmt.Printf("Logout failed: %v", err)

			// Handle specific errors
			statusCode := http.StatusInternalServerError
			if errors.Is(err, errors.New("refresh token is empty")) {
				statusCode = http.StatusNotFound
			}
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    statusCode,
				Message: "Failed to log out",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusBadRequest,
				Message: "Logged out successfully",
			},
			Data: nil,
		})

	}
}
