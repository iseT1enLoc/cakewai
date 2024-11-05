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

func NewRoleRoute(env *appconfig.Env, timeout time.Duration, db *mongo.Database, r *gin.RouterGroup) {
	repo := repository.NewRoleRepository(db, "roles")
	sc := handlers.RoleController{
		RoleUsecase: usecase.NewRoleUsecase(repo, timeout),
		Env:         env,
	}
	r.POST("/role", sc.CreateRole())              //ok
	r.GET("/roles", sc.GetAllRoles())             //ok
	r.GET("/role/:id", sc.GetRoleByID())          //ok
	r.GET("/rolex/:name", sc.GetRoleByRoleName()) //ok
	r.PUT("/role/:id", sc.UpdateRole())
}
