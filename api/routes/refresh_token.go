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

func NewRefreshTokenRoute(env *appconfig.Env, timeout time.Duration, db *sql.DB, r *gin.RouterGroup) {
	repo := repository.NewUserRepository(db)
	sc := handlers.RefreshTokenHandler{
		RefreshTokenUsecase: usecase.NewRefreshTokenUseCase(repo, timeout),
		Env:                 env,
	}
	r.POST("/refreshtoken", sc.RefreshTokenHandler())
}
