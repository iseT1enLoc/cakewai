package routes

import (
	"cakewai/cakewai.com/api/handlers"
	"cakewai/cakewai.com/api/middlewares"
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/repository"
	"cakewai/cakewai.com/usecase"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

func NewRefreshTokenRoute(env *appconfig.Env, timeout time.Duration, db *mongo.Database, r *gin.RouterGroup) {
	repo := repository.NewUserRepository(db, "users")
	sc := handlers.RefreshTokenHandler{
		RefreshTokenUsecase: usecase.NewRefreshTokenUseCase(repo, timeout),
		Env:                 env,
	}
	r.GET("/refreshtoken", middlewares.TraceMiddleware("refresh token middlware"), sc.RefreshTokenHandler())
}
