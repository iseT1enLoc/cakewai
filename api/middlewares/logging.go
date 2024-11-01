package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

func TraceMiddleware(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Entering %s", name) // Log entering a handler or middleware

		c.Next() // Move to the next handler or middleware

		log.Printf("Exiting %s", name) // Log exiting after handler or middleware completes
	}
}
