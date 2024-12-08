package routes

import (
	"cakewai/cakewai.com/api/handlers"
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/repository"
	"cakewai/cakewai.com/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewLogoutRoute(env *appconfig.Env, timeout time.Duration, db *mongo.Database, r *gin.RouterGroup) {
	RT_repository := repository.NewrefreshTokenRepository(db, "refresh_token")
	sc := handlers.LogoutHandler{
		LogoutUsecase: usecase.NewLogoutUseCase(RT_repository, timeout),
		Env:           env,
	}
	r.POST("/logout", sc.LogoutHandler())
}
