package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/domain"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RefreshTokenHandler struct {
	RefreshTokenUsecase domain.RefreshTokenUseCase
	Env                 *appconfig.Env
}
type onlyRRefreshRequest struct {
	TokenID string `json:"refresh_token" bson:"refresh_token" form:"refresh_token" binding:"required"`
}

func (u *RefreshTokenHandler) RefreshTokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		var reqToken onlyRRefreshRequest
		//Bind JSON request body to RefreshTokenRequest
		if err := c.ShouldBindJSON(&reqToken); err != nil {

			fmt.Println("Error at line 23")
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Can not parsing from json",
				Error:   err.Error(),
			})
			return
		}
		refresh_token, err := u.RefreshTokenUsecase.GetRefreshTokenFromDB(c, reqToken.TokenID, u.Env)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid refresh token",
				Error:   err.Error(),
			})
			return
		}

		if time.Now().Local().After(refresh_token.ExpireAt.Local()) || refresh_token.IsExpire {

			accessToken, refresh_token, err := u.RefreshTokenUsecase.RefreshToken(c, *refresh_token, u.Env)

			fmt.Println(accessToken)
			fmt.Println(refresh_token)
			if err != nil {
				c.JSON(http.StatusBadRequest, response.FailedResponse{
					Code:    http.StatusBadRequest,
					Message: "Error happened after refresh token",
					Error:   err.Error(),
				})
				return
			}

			responsetoken := domain.RefreshShortResponse{
				AccessToken:  accessToken,
				RefreshToken: refresh_token,
			}

			c.JSON(http.StatusOK, response.Success{
				ResponseFormat: response.ResponseFormat{
					Code:    200,
					Message: "Refresh token",
				},
				Data: responsetoken,
			})
			return

		}

		accessToken, refreshtoken, err := u.RefreshTokenUsecase.RenewAccessToken(c, *refresh_token, u.Env)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error happened while renew access token",
				Error:   err.Error(),
			})
			return
		}
		responsetoken := domain.RefreshShortResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshtoken,
		}
		c.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully renew token",
			},
			Data: responsetoken,
		})

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
