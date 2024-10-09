package routes

import (
	"cakewai/cakewai.com/api/handlers"
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/repository"
	"cakewai/cakewai.com/usecase"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(env *appconfig.Env, timeout time.Duration, db *sql.DB, r *gin.RouterGroup) {
	ur := repository.NewUserRepository(db)
	// uc := &controller.UserController{
	// 	UserUseCase: usecase.NewUserUseCase(ur, timeout),
	// 	Env:         env,
	// }
	uc := &handlers.UserController{
		UserUseCase: usecase.NewProfileUseCase(timeout, ur),
		Env:         env,
	}
	// USER ROUTES
	group := r.Group("/user")
	{
		group.GET("/profile", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "congratulation for joining our platform"})
		})
		group.GET("/all", uc.GetUsers())           // Get all users
		group.GET("/:user_id", uc.GetUserById())   // Get user by ID
		group.PUT("/:user_id", uc.UpdateUser())    // Update user by ID
		group.DELETE("/:user_id", uc.DeleteUser()) // Delete user by ID
	}
}
