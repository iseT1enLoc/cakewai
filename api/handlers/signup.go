package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/domain"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ydb-platform/ydb-go-sdk/v3/log"
)

type SignupController struct {
	SignupUseCase domain.SignupUseCase
	Env           *appconfig.Env
}

func (sc *SignupController) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(c, time.Second*time.Duration(c.GetFloat64(sc.Env.TIMEOUT)))
		defer cancel()
		var request domain.SignupRequest
		// Decode the JSON request body
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Can not parsing the request",
				Error:   err.Error(),
			})
			return
		}

		// Call the signup use case
		accessToken, refreshToken, err := sc.SignupUseCase.SignUp(c, request, sc.Env)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Email is already existed",
				Error:   err.Error(),
			})
			return
		}
		// Create the response object
		signupResponse := domain.SignupResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}
		// Send the response
		c.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusCreated,
				Message: "Register successfully",
			},
			Data: signupResponse,
		})
	}
}
