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

func NewAccountRoute(Env *appconfig.Env, timout time.Duration, db *mongo.Database, r *gin.RouterGroup) {
	user_repo := repository.NewUserRepository(db, "users")
	acc_recover_handler := handlers.AccountRecoveryHandler{
		Acc_recover_usecase: usecase.NewAccountRecovery(user_repo, timout),
		Env:                 Env,
	}

	r.POST("/request-password-reset", acc_recover_handler.ResetPasswordRequest())
	r.POST("/reset-password", acc_recover_handler.ResetPasswordProcessing())
}
