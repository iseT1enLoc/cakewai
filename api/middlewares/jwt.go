package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"cakewai/cakewai.com/component/response"
	tokenutil "cakewai/cakewai.com/internals/token_utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func extractClaim(authToken string, secret []byte) (string, error) {
	// Parse the token
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return "", err
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Assuming "userID" is the claim you want to extract
		userID, ok := claims["userID"].(string)
		if !ok {
			return "", fmt.Errorf("userID claim is not found or not a string")
		}
		return userID, nil
	}
	return "", fmt.Errorf("invalid token or unable to extract claims")
}

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
				userID, is_admin, err := tokenutil.ExtractIDAndRole(authToken, secret)

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
				// Set user ID in context
				c.Set("is_admin", is_admin)
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

// AdminMiddleware checks if the user is an admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the user ID from the context set in JwtAuthMiddleware
		is_admin, exists := c.Get("is_admin")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "is_admin not found"})
			c.Abort()
			return
		}

		// Here, implement your logic to check if the user is an admin
		// This could involve a DB check or a predefined list of admin IDs.
		// Example: checking a list of admin userIDs
		if is_admin == false {
			c.JSON(http.StatusForbidden, gin.H{"message": "Admin access required"})
			c.Abort()
			return
		}

		// If the user is an admin, continue to the next handler
		c.Next()
	}
}

// AdminMiddleware checks if the user is an admin
func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the user ID from the context set in JwtAuthMiddleware
		is_admin, exists := c.Get("is_admin")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "is_admin not found"})
			c.Abort()
			return
		}

		// Here, implement your logic to check if the user is an admin
		// This could involve a DB check or a predefined list of admin IDs.
		// Example: checking a list of admin userIDs
		if is_admin == true {
			c.JSON(http.StatusForbidden, gin.H{"message": "User route"})
			c.Abort()
			return
		}

		// If the user is an admin, continue to the next handler
		c.Next()
	}
}
