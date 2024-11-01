package routes

import (
	"cakewai/cakewai.com/api/handlers"
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/repository"
	"cakewai/cakewai.com/usecase"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

func NewLoginRoute(env *appconfig.Env, timeout time.Duration, db *mongo.Database, r *gin.RouterGroup) {
	repo := repository.NewUserRepository(db, "users")
	RT_repository := repository.NewrefreshTokenRepository(db, "refresh_token")
	sc := handlers.LoginHandler{
		LoginUsecase: usecase.NewLoginUseCase(repo, RT_repository, timeout),
		Env:          env,
	}
	r.POST("/login", sc.LoginHandler())
}
