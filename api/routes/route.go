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

	//protectedRoute.Use(middlewares.JwtAuthMiddleware(os.Getenv("ACCESS_SECRET")))
	NewGoogleRouter(env, timeout, db, publicRoute)
	NewSignUpRoute(env, timeout, db, publicRoute)
	NewLoginRoute(env, timeout, db, publicRoute)
	protectedRoute.Use(middlewares.JwtAuthMiddleware("hihi"))

	NewUserRouter(env, timeout, db, protectedRoute)
	NewRefreshTokenRoute(env, timeout, db, publicRoute)
	NewProductRoute(env, timeout, db, publicRoute)
	//NewSignInRoute(env, timeout, db, publicRoute)

	//NewResourceRoute(env, timeout, db, protectedRoute)

}
