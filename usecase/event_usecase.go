package usecase

import (
	"context"
	"errors"
	"time"

	"cakewai/cakewai.com/domain"
	"cakewai/cakewai.com/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventBlogUsecase interface defines the business operations for EventBlog
type EventBlogUsecase interface {
	CreateEventBlog(ctx context.Context, eventBlog domain.EventBlog) (*domain.EventBlog, error)
	GetEventBlogById(ctx context.Context, id primitive.ObjectID) (*domain.EventBlog, error)
	GetAllEventBlogs(ctx context.Context) ([]domain.EventBlog, error)
	UpdateEventBlog(ctx context.Context, eventBlog domain.EventBlog) (*domain.EventBlog, error)
	DeleteEventBlog(ctx context.Context, id primitive.ObjectID) error
}

// eventBlogUsecase struct implements the EventBlogUsecase interface
type eventBlogUsecase struct {
	eventBlogRepo repository.EventBlogRepository
}

// NewEventBlogUsecase creates a new instance of eventBlogUsecase
func NewEventBlogUsecase(repo repository.EventBlogRepository) EventBlogUsecase {
	return &eventBlogUsecase{
		eventBlogRepo: repo,
	}
}

// CreateEventBlog creates a new event blog
func (uc *eventBlogUsecase) CreateEventBlog(ctx context.Context, eventBlog domain.EventBlog) (*domain.EventBlog, error) {
	if eventBlog.Title == "" {
		return nil, errors.New("title is required")
	}
	if eventBlog.ShortDescription == "" {
		return nil, errors.New("short description is required")
	}

	// Set the created_at field if not set
	if eventBlog.CreatedAt.IsZero() {
		eventBlog.CreatedAt = time.Now()
	}

	// Call the repository to insert the event blog into the database
	return uc.eventBlogRepo.CreateEventBlog(ctx, eventBlog)
}

// GetEventBlogById retrieves an event blog by its ID
func (uc *eventBlogUsecase) GetEventBlogById(ctx context.Context, id primitive.ObjectID) (*domain.EventBlog, error) {
	// Call the repository to retrieve the event blog by ID
	eventBlog, err := uc.eventBlogRepo.GetEventBlogById(ctx, id)
	if err != nil {
		return nil, err
	}
	return eventBlog, nil
}

// GetAllEventBlogs retrieves all event blogs
func (uc *eventBlogUsecase) GetAllEventBlogs(ctx context.Context) ([]domain.EventBlog, error) {
	// Call the repository to retrieve all event blogs
	eventBlogs, err := uc.eventBlogRepo.GetAllEventBlogs(ctx)
	if err != nil {
		return nil, err
	}
	return eventBlogs, nil
}

// UpdateEventBlog updates an existing event blog
func (uc *eventBlogUsecase) UpdateEventBlog(ctx context.Context, eventBlog domain.EventBlog) (*domain.EventBlog, error) {
	// Ensure the event blog ID is set before updating
	if eventBlog.Id.IsZero() {
		return nil, errors.New("event blog ID is required")
	}

	// Call the repository to update the event blog
	return uc.eventBlogRepo.UpdateEventBlog(ctx, eventBlog)
}

// DeleteEventBlog deletes an event blog by ID
func (uc *eventBlogUsecase) DeleteEventBlog(ctx context.Context, id primitive.ObjectID) error {
	// Call the repository to delete the event blog by ID
	err := uc.eventBlogRepo.DeleteEventBlog(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
