package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventBlog struct {
	Id               primitive.ObjectID `json:"_id" bson:"_id"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	ImageLink        string             `json:"image_link" bson:"image_link"`
	Title            string             `json:"title" bson:"title"`
	Author           string             `json:"author" bson:"author"`
	ShortDescription string             `json:"short_description" bson:"short_description"`
}

type EventUsecase interface {
	CreateEventBlog(ctx context.Context, title, shortDescription string) (*EventBlog, error)
	GetEventBlogById(ctx context.Context, id string) (*EventBlog, error)
	GetAllEventBlogs(ctx context.Context) ([]EventBlog, error)
}
