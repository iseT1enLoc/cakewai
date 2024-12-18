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

func NewEventBlogRoute(Env *appconfig.Env, timeout time.Duration, db *mongo.Database, r *gin.RouterGroup) {
	// Create repository and handler for Event Blog
	event_blog_repo := repository.NewEventBlogRepository(db, "event_blog")
	event_blog_handler := handlers.EventBlogHandler{
		EventBlogUsecase: usecase.NewEventBlogUsecase(event_blog_repo),
	}

	// Create EventBlog
	r.POST("/event_blog", event_blog_handler.CreateEventBlogHandler()) // Create a new event blog

	// Get EventBlog by ID
	r.GET("/event_blog/:id", event_blog_handler.GetEventBlogByIdHandler()) // Get event blog by ID

	// Get all EventBlogs
	r.GET("/event_blog", event_blog_handler.GetAllEventBlogsHandler()) // Get all event blogs

	// Update EventBlog
	r.PUT("/event_blog/:id", event_blog_handler.UpdateEventBlogHandler()) // Update an event blog by ID

	// Delete EventBlog
	r.DELETE("/event_blog/:id", event_blog_handler.DeleteEventBlogHandler()) // Delete an event blog by ID
}
