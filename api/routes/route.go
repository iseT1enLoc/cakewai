package routes

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

func SetUp(env *appconfig.Env, timeout time.Duration, db *mongo.Database, r *gin.Engine) {
	publicRoute := r.Group("/api/public")
	protectedRoute := r.Group("/api/protected")

	//protectedRoute.Use(middlewares.JwtAuthMiddleware(os.Getenv("ACCESS_SECRET")))
	NewGoogleRouter(env, timeout, db, publicRoute)
	NewSignUpRoute(env, timeout, db, publicRoute)
	NewLoginRoute(env, timeout, db, publicRoute)
	NewUserRouter(env, timeout, db, protectedRoute)
	NewRefreshTokenRoute(env, timeout, db, publicRoute)
	//NewSignInRoute(env, timeout, db, publicRoute)
	//protectedRoute.Use(middleware.JwtAuthMiddleware(os.Getenv("SECRET_KEY")))
	//NewResourceRoute(env, timeout, db, protectedRoute)
}
