package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"cakewai/cakewai.com/component/response"
	tokenutil "cakewai/cakewai.com/internals/token_utils"

	"github.com/gin-gonic/gin"
)

// JwtAuthMiddleware is the Gin middleware for JWT authentication
func JwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println("Enter line 17 jwt auth middleware")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			print(authToken)
			// Validate the token
			authorized, err := tokenutil.Is_authorized(authToken, secret)
			if err != nil {
				fmt.Printf("Print line 24 at jwt middleware")
				c.JSON(http.StatusUnauthorized, response.FailedResponse{
					Code:    http.StatusUnauthorized,
					Message: "Is not authorized",
					Error:   err.Error(),
				})
				fmt.Println(err)
				c.Abort()
				return
			}

			fmt.Printf("print at line 29 jwtAuthMiddleWare")
			if authorized {
				fmt.Println("authorized")
				// Extract user ID from token
				userID, err := tokenutil.ExtractID(authToken, secret)

				fmt.Printf("print user id after extractID %v", userID)
				if err != nil {
					c.JSON(http.StatusForbidden, response.FailedResponse{
						Code:    http.StatusForbidden,
						Message: "Fail to extract id",
						Error:   err.Error(),
					})
					c.Abort()
					return
				}

				// Set user ID in context
				c.Set("user_id", userID)

				fmt.Printf("USER ID: %s\n", userID)
				c.Next() // Continue to the next handler
				fmt.Println("After c.next")
				return
			}
			fmt.Printf("print at line 44 jwtAuthMiddleWare")
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthrized"})
			c.Abort()
			return
		}

		// If Authorization header is missing or malformed
		c.JSON(http.StatusUnauthorized, gin.H{"message": "malform existed"})
		c.Abort()
	}
}
