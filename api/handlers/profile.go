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
		fmt.Printf("User id is that %v", userID)
		if !exists {
			fmt.Print("User ID not found in context")
			return
		}

		currentuser, err := uc.UserUseCase.GetUserById(c, userID.(string))
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
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    0,
				Message: "Error while query database",
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully get all users",
			},
			Data: users,
		})
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
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while query database",
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully get user by id",
			},
			Data: user,
		})
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
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Can not convert user id to hexa",
				Error:   err.Error(),
			})
			return
		}

		// Decode the user from the request body
		var user domain.User
		if err := c.ShouldBindJSON(&user); err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Can not parsing",
				Error:   err.Error(),
			})
			return
		}
		fmt.Println(user)
		// Set the user ID
		user.Id = userId

		// Call the use case to update the user
		if err := uc.UserUseCase.UpdateUser(c, &user); err != nil {
			log.Error(err)
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while querying database",
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Update user successfully",
			},
			Data: user.Id,
		})
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
			c.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Error while querying database",
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully delete the user",
			},
			Data: userIDParam,
		})
	}
}
