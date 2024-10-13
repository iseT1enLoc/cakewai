package routes

import (
	"cakewai/cakewai.com/api/handlers"
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/repository"
	"cakewai/cakewai.com/usecase"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

func NewSignUpRoute(env *appconfig.Env, timeout time.Duration, db *mongo.Database, r *gin.RouterGroup) {
	fmt.Print("singup")
	repo := repository.NewUserRepository(db, "users")
	sc := handlers.SignupController{
		SignupUseCase: usecase.NewSignupUseCase(repo, timeout),
		Env:           env,
	}
	r.POST("/signup", sc.SignUp())
}
