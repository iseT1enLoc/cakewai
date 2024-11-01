package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/domain"
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"github.com/ydb-platform/ydb-go-sdk/v3/log"
)

type UserController struct {
	UserUseCase domain.UserUseCase
	Env         *appconfig.Env
}

func (uc *UserController) GetCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		//user_id := c.GetString("user_id")
		userID, exists := c.Get("user_id")
		if !exists {
			fmt.Print("User ID not found in context")
			return
		}
		userIDStr := userID.(string) // Type assertion if it’s a string
		fmt.Printf("/n The current user is %v", userIDStr)
		fmt.Println()
		currentuser, err := uc.UserUseCase.GetUserById(c, "672167f95b55a7c71f00f18a")
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    401,
				Message: "Invalid user id",
				Error:   "can not get user from database",
			})
			return
		}
		c.JSON(http.StatusOK, response.SuccessResponse{
			Success: response.Success{
				ResponseFormat: response.ResponseFormat{
					Code:    http.StatusOK,
					Message: "this is your user",
				},
				Data: currentuser,
			},
			Meta: response.Meta{},
		})
	}
}

func (uc *UserController) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from the context
		userID := c.Param("user_id") // Assuming user_id is part of the URL parameters

		// Optionally, you can also set it in the context
		ctx := context.WithValue(c.Request.Context(), "user_id", userID)

		users, err := uc.UserUseCase.GetListUsers(ctx)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}
func (uc *UserController) GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from the URL parameters (assuming it's part of the path)
		sdParam := c.Param("user_id")
		print(sdParam)
		fmt.Print("get user by id handler line 42")
		// Convert the user ID from string to int
		// intId, err := uuid.FromBytes([]byte(idParam))
		// if err != nil {
		// 	log.Error(err)
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		// 	return
		// }
		fmt.Print("get user by id line 50")
		// Call the use case to get the user by ID
		user, err := uc.UserUseCase.GetUserById(c.Request.Context(), sdParam)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
func (uc *UserController) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from the context
		userIDParam := c.Param("user_id") // Assuming user_id is part of the URL parameters

		// Convert user ID from string to int
		userId, err := primitive.ObjectIDFromHex(userIDParam)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
			return
		}

		// Decode the user from the request body
		var user domain.User
		if err := c.ShouldBindJSON(&user); err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set the user ID
		user.Id = userId

		// Call the use case to update the user
		if err := uc.UserUseCase.UpdateUser(c.Request.Context(), &user); err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	}
}
func (uc *UserController) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from the URL parameters
		userIDParam := c.Param("user_id")

		// Convert user ID from string to int
		// userId, err := strconv.Atoi(userIDParam)
		// if err != nil {
		// 	log.Error(err)
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		// 	return
		// }

		// Call the use case to delete the user
		if err := uc.UserUseCase.DeleteUser(c.Request.Context(), userIDParam); err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	}
}
