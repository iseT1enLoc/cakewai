package handlers

import (
	"cakewai/cakewai.com/domain"
	"cakewai/cakewai.com/usecase"

	"github.com/gin-gonic/gin"

	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Success represents a successful response structure
type Success struct {
	ResponseFormat
	Data interface{} `json:"data,omitempty"`
}

// FailedResponse represents a failed response structure
type FailedResponse struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"internal_server_error"`
	Error   string `json:"error" example:"{$err}"`
}

// ResponseFormat is the base format for all responses
type ResponseFormat struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"success"`
}

// EventBlogHandler struct to bind the usecase with the Gin handler
type EventBlogHandler struct {
	EventBlogUsecase usecase.EventBlogUsecase
}

// NewEventBlogHandler creates a new EventBlogHandler
func NewEventBlogHandler(usecase usecase.EventBlogUsecase) *EventBlogHandler {
	return &EventBlogHandler{
		EventBlogUsecase: usecase,
	}
}

// CreateEventBlogHandler handles the creation of an EventBlog
func (h *EventBlogHandler) CreateEventBlogHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var eventBlog domain.EventBlog

		if err := c.ShouldBindJSON(&eventBlog); err != nil {
			c.JSON(http.StatusBadRequest, FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "bad_request",
				Error:   err.Error(),
			})
			return
		}

		createdEventBlog, err := h.EventBlogUsecase.CreateEventBlog(c, eventBlog)
		if err != nil {
			log.Println("Error creating event blog:", err)
			c.JSON(http.StatusInternalServerError, FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "internal_server_error",
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, Success{
			ResponseFormat: ResponseFormat{
				Status:  http.StatusOK,
				Message: "success",
			},
			Data: createdEventBlog,
		})
	}

}

// GetEventBlogByIdHandler handles retrieving an EventBlog by its ID
func (h *EventBlogHandler) GetEventBlogByIdHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid_id_format",
				Error:   err.Error(),
			})
			return
		}

		eventBlog, err := h.EventBlogUsecase.GetEventBlogById(c, id)
		if err != nil {
			log.Println("Error retrieving event blog:", err)
			c.JSON(http.StatusInternalServerError, FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "internal_server_error",
				Error:   err.Error(),
			})
			return
		}

		if eventBlog == nil {
			c.JSON(http.StatusNotFound, FailedResponse{
				Code:    http.StatusNotFound,
				Message: "event_blog_not_found",
				Error:   "event blog with the given ID was not found",
			})
			return
		}

		c.JSON(http.StatusOK, Success{
			ResponseFormat: ResponseFormat{
				Status:  http.StatusOK,
				Message: "success",
			},
			Data: eventBlog,
		})
	}

}

// GetAllEventBlogsHandler handles retrieving all EventBlogs
func (h *EventBlogHandler) GetAllEventBlogsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		eventBlogs, err := h.EventBlogUsecase.GetAllEventBlogs(c)
		if err != nil {
			log.Println("Error retrieving event blogs:", err)
			c.JSON(http.StatusInternalServerError, FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "internal_server_error",
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, Success{
			ResponseFormat: ResponseFormat{
				Status:  http.StatusOK,
				Message: "success",
			},
			Data: eventBlogs,
		})
	}

}

// UpdateEventBlogHandler handles updating an existing EventBlog
func (h *EventBlogHandler) UpdateEventBlogHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var eventBlog domain.EventBlog

		if err := c.ShouldBindJSON(&eventBlog); err != nil {
			c.JSON(http.StatusBadRequest, FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "bad_request",
				Error:   err.Error(),
			})
			return
		}

		eventBlogId := c.Param("id")
		id, err := primitive.ObjectIDFromHex(eventBlogId)
		if err != nil {
			c.JSON(http.StatusBadRequest, FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid_id_format",
				Error:   err.Error(),
			})
			return
		}

		eventBlog.Id = id

		updatedEventBlog, err := h.EventBlogUsecase.UpdateEventBlog(c, eventBlog)
		if err != nil {
			log.Println("Error updating event blog:", err)
			c.JSON(http.StatusInternalServerError, FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "internal_server_error",
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, Success{
			ResponseFormat: ResponseFormat{
				Status:  http.StatusOK,
				Message: "success",
			},
			Data: updatedEventBlog,
		})
	}

}

// DeleteEventBlogHandler handles deleting an EventBlog by its ID
func (h *EventBlogHandler) DeleteEventBlogHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid_id_format",
				Error:   err.Error(),
			})
			return
		}

		err = h.EventBlogUsecase.DeleteEventBlog(c, id)
		if err != nil {
			log.Println("Error deleting event blog:", err)
			c.JSON(http.StatusInternalServerError, FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "internal_server_error",
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, Success{
			ResponseFormat: ResponseFormat{
				Status:  http.StatusOK,
				Message: "success",
			},
			Data: nil,
		})
	}

}
