package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	tokenutil "cakewai/cakewai.com/internals/token_utils"

	"github.com/gin-gonic/gin"
)

// JwtAuthMiddleware is the Gin middleware for JWT authentication
func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			print(authToken)
			// Validate the token
			authorized, err := tokenutil.IsAuthorized(authToken, secret)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"Error": "Validate token"})
				c.Abort()
				return
			}

			if authorized {
				// Extract user ID from token
				userID, err := tokenutil.ExtractID(authToken, secret)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"Error": "extract id"})
					c.Abort()
					return
				}

				// Set user ID in context
				c.Set("user_id", userID)
				fmt.Printf("USER ID: %s\n", userID)
				c.Next() // Continue to the next handler
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthrized"})
			c.Abort()
			return
		}

		// If Authorization header is missing or malformed
		c.JSON(http.StatusUnauthorized, gin.H{"message": "malform existed"})
		c.Abort()
	}
}
