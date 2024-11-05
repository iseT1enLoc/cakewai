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

func NewSignUpRoute(env *appconfig.Env, timeout time.Duration, db *mongo.Database, r *gin.RouterGroup) {
	repo := repository.NewUserRepository(db, "users")
	repo_refresh := repository.NewrefreshTokenRepository(db, "refresh_token")
	sc := handlers.SignupController{
		SignupUseCase: usecase.NewSignupUseCase(repo, repo_refresh, timeout),
		Env:           env,
	}
	r.POST("/signup", sc.SignUp())
}
