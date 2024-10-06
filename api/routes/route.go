package routes

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
)

func SetUp(env *appconfig.Env, timeout time.Duration, db *sql.DB, r *gin.Engine) {
	publicRoute := r.Group("/api/public/")
	//protectedRoute := r.Group("api/protected/")

	NewGoogleRouter(env, timeout, db, publicRoute)
	NewSignUpRoute(env, timeout, db, publicRoute)
	//NewSignInRoute(env, timeout, db, publicRoute)
	//protectedRoute.Use(middleware.JwtAuthMiddleware(os.Getenv("SECRET_KEY")))
	//NewResourceRoute(env, timeout, db, protectedRoute)
}
