package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/domain"
	"errors"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

type AccountRecoveryHandler struct {
	Acc_recover_usecase domain.AccountRecovery
	Env                 *appconfig.Env
}
type NewPassword struct {
	NewPassword string `json:"new_password" bson:"new_password"`
	Token       string `json:"token" bson:"token"`
}

func (a *AccountRecoveryHandler) ResetPasswordProcessing() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("Enter reset password processing line 25")
		var new_password NewPassword
		if err := ctx.ShouldBind(&new_password); err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request body",
				Error:   err.Error(),
			})
			return
		}
		fmt.Println("Enter reset password processing line 355")
		log.Println(new_password)
		fmt.Println("Enter reset password processing line 25")
		err := a.Acc_recover_usecase.ResetPasswordProcessing(ctx, a.Env, new_password.NewPassword, new_password.Token)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Password reset processing fail",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Password is updated successfully, login again with your new password",
			},
			Data: nil,
		})
	}
}
func (a *AccountRecoveryHandler) ResetPasswordRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req domain.ResetPasswordReq
		//parsing to get the body
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request body",
				Error:   err.Error(),
			})
			return
		}

		err := a.Acc_recover_usecase.ResetPasswordRequest(ctx, a.Env, req.Email)
		//parsing to get the body
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while handling request",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Email is sent to user mail box",
			},
			Data: nil,
		})
	}
}
func (a *AccountRecoveryHandler) ChangesPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("user_id")
		fmt.Printf("User id is that %v", userID)
		if !exists {
			fmt.Print("User ID not found in context")
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid user param",
				Error:   errors.New("Can not find user id in context").Error(),
			})
			return
		}
		objhex, err := primitive.ObjectIDFromHex(userID.(string))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid user param",
				Error:   err.Error(),
			})
		}

		type currentPassword struct {
			CurrentPassword string `json:"current_password" bson:"password"`
			NewPassword     string `json:"new_password" bson:"new_password"`
		}
		var existingPassword currentPassword
		if err := ctx.ShouldBind(&existingPassword); err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request body",
				Error:   err.Error(),
			})
			return
		}

		fmt.Println(existingPassword.CurrentPassword)
		fmt.Println(existingPassword.NewPassword)
		err = a.Acc_recover_usecase.ChangesPassword(ctx, a.Env, objhex, existingPassword.CurrentPassword, existingPassword.NewPassword)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Can not change password",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully change password",
			},
			Data: nil,
		})
	}
}
