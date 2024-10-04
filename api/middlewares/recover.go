package middlewares

import (
	appctx "cakewai/cakewai.com/component/appcontext"
	apperror "cakewai/cakewai.com/component/apperr"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recover(appCtx appctx.Appcontext) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")

				if appErr, ok := err.(*apperror.AppError); ok {
					c.AbortWithStatusJSON(appErr.Code, appErr)
					panic(err)
				}

				appErr := apperror.InternalServerError(err.(error))

				c.AbortWithStatusJSON(http.StatusInternalServerError, appErr)
				panic(err)
			}
		}()

		c.Next()
	}
}
