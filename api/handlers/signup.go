package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/domain"
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"github.com/ydb-platform/ydb-go-sdk/v3/log"
)

type SignupController struct {
	SignupUseCase domain.SignupUseCase
	CartUseCase   domain.CartUsecase
	Env           *appconfig.Env
}

func (sc *SignupController) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Context timeout handling
		ctx, cancel := context.WithTimeout(c, time.Second*5)
		defer cancel()

		// Decode the JSON request body
		var request domain.SignupRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Cannot parse the request body",
				Error:   err.Error(),
			})
			return
		}

		// Set default role if it's empty
		if request.RoleName == "" {
			request.RoleName = "customer"
		}

		// Call the signup use case
		accessToken, refreshToken, uid, err := sc.SignupUseCase.SignUp(ctx, request, sc.Env)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Email already exists",
				Error:   err.Error(),
			})
			return
		}

		// Create the response object
		signupResponse := domain.SignupResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		// Convert the user ID from string to ObjectID and create an empty cart for the user
		hexid, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid user ID",
				Error:   err.Error(),
			})
			return
		}

		err = sc.CartUseCase.CreateCartByUserId(ctx, hexid)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error creating empty cart",
				Error:   err.Error(),
			})
			return
		}

		// Send the success response
		c.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusCreated,
				Message: "Registered successfully",
			},
			Data: signupResponse,
		})
	}
}
