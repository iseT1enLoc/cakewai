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

func NewGoogleRouter(env *appconfig.Env, timeout time.Duration, db *mongo.Database, r *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, "users")
	gc := &handlers.GoogleController{
		GoogleUseCase: usecase.NewGoogleUseCase(ur, timeout),
		Env:           env,
	}

	r.GET("/google/login", gc.HandleGoogleLogin())
	r.GET("/google/callback", gc.HandleGoogleCallback())

}
