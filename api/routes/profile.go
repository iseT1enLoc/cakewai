package routes

import (
	"cakewai/cakewai.com/api/handlers"
	"cakewai/cakewai.com/api/middlewares"
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/repository"
	"cakewai/cakewai.com/usecase"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(env *appconfig.Env, timeout time.Duration, db *mongo.Database, r *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, "users")
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
		group.GET("/current_user", middlewares.TraceMiddleware("profle"), uc.GetCurrentUser())
		group.GET("/all", middlewares.TraceMiddleware("get all user middleware"), uc.GetUsers())          // Get all users
		group.GET("/:user_id", middlewares.TraceMiddleware("get user id middle ware"), uc.GetUserById())  // Get user by ID
		group.PUT("/:user_id", middlewares.TraceMiddleware("user id update middleware"), uc.UpdateUser()) // Update user by ID
		group.DELETE("/:user_id", middlewares.TraceMiddleware("delete user middleware"), uc.DeleteUser()) // Delete user by ID
	}
}
