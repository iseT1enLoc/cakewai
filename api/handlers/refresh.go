package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/domain"
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

func (u *RefreshTokenHandler) RenewRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1: Parse and validate the request body
		var reqToken onlyRRefreshRequest
		if err := c.ShouldBindJSON(&reqToken); err != nil {
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid JSON input",
				Error:   err.Error(),
			})
			return
		}

		// // Create an instance of the custom claims struct
		// claims := &domain.JwtCustomClaims{}

		// // Parse the token string with claims
		// _, err := jwt.ParseWithClaims(req.TokenID, claims, func(token *jwt.Token) (interface{}, error) {
		// 	// Validate the signing method
		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		// 	}
		// 	// Return the secret key
		// 	return []byte(u.Env.REFRESH_SECRET), nil
		// })
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, response.FailedResponse{
		// 		Code:    http.StatusBadRequest,
		// 		Message: "Error while persing token",
		// 		Error:   err.Error(),
		// 	})
		// 	return
		// }
		accessToken, refreshToken, err := u.RefreshTokenUsecase.RefreshToken(c.Request.Context(), reqToken.TokenID, false, u.Env)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while creating new refresh token ",
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Refresh token renewed successfully",
			},
			Data: domain.RefreshShortResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			},
		})
	}
}
func (u *RefreshTokenHandler) RenewAcessTokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1: Parse and validate the request body
		var reqToken onlyRRefreshRequest
		if err := c.ShouldBindJSON(&reqToken); err != nil {
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid JSON input",
				Error:   err.Error(),
			})
			return
		}

		// Step 2: Retrieve the refresh token from the database
		refreshToken, err := u.RefreshTokenUsecase.GetRefreshTokenFromDB(c, reqToken.TokenID, u.Env)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid refresh token",
				Error:   err.Error(),
			})
			return
		}

		// Step 3: Validate the refresh token's expiration status
		if time.Now().After(refreshToken.ExpireAt) || refreshToken.IsExpire {
			c.JSON(http.StatusUnauthorized, response.FailedResponse{
				Code:    http.StatusUnauthorized,
				Message: "Refresh token expired",
				Error:   "Token is no longer valid",
			})
			return
		}

		// Step 4: Renew access and refresh tokens
		accessToken, newRefreshToken, err := u.RefreshTokenUsecase.RenewAccessToken(c, *refreshToken, u.Env)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to renew access token",
				Error:   err.Error(),
			})
			return
		}

		// Step 5: Respond with new tokens
		c.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Access Tokens renewed successfully",
			},
			Data: domain.RefreshShortResponse{
				AccessToken:  accessToken,
				RefreshToken: newRefreshToken,
			},
		})
	}
}
func (u *RefreshTokenHandler) RevokeRefreshTokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request domain.RefreshTokenRequest

		// Parse the JSON request body
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    0,
				Message: "Invalid request payload",
				Error:   err.Error(),
			})
			return
		}

		// Attempt to revoke the refresh token
		if err := u.RefreshTokenUsecase.RevokeToken(c, request.RefreshToken, u.Env); err != nil {
			c.JSON(http.StatusUnauthorized, response.FailedResponse{
				Code:    0,
				Message: "Failed to revoke refresh token",
				Error:   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Refresh token successfully revoked",
			},
			Data: gin.H{"status": "success"},
		})

	}
}
