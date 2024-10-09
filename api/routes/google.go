package routes

import (
	"cakewai/cakewai.com/api/handlers"
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/repository"
	"cakewai/cakewai.com/usecase"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
)

func NewGoogleRouter(env *appconfig.Env, timeout time.Duration, db *sql.DB, r *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	gc := &handlers.GoogleController{
		GoogleUseCase: usecase.NewGoogleUseCase(ur, timeout),
		Env:           env,
	}

	r.GET("/google/login", gc.HandleGoogleLogin())
	r.GET("/google/callback", gc.HandleGoogleCallback())

}
