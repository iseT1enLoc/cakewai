package routes

import (
	"cakewai/cakewai.com/api/middlewares"
	appconfig "cakewai/cakewai.com/component/appcfg"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

func SetUp(env *appconfig.Env, timeout time.Duration, db *mongo.Database, r *gin.Engine) {
	publicRoute := r.Group("/api/public")       //ai vo cung duoc
	protectedRoute := r.Group("/api/protected") // co account moi vo duoc

	NewGoogleRouter(env, timeout, db, publicRoute)
	NewSignUpRoute(env, timeout, db, publicRoute)
	NewLoginRoute(env, timeout, db, publicRoute)
	protectedRoute.Use(middlewares.JwtAuthMiddleware(env.ACCESS_SECRET))

	NewUserRouter(env, timeout, db, protectedRoute)
	NewRefreshTokenRoute(env, timeout, db, publicRoute)
	NewProductRoute(env, timeout, db, publicRoute)
	NewRoleRoute(env, timeout, db, publicRoute)
	NewCartRoute(env, timeout, db, protectedRoute)
}
