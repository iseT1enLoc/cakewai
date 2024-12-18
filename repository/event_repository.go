package repository

import (
	"context"

	"cakewai/cakewai.com/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventBlogRepository interface {
	CreateEventBlog(ctx context.Context, eventBlog domain.EventBlog) (*domain.EventBlog, error)
	GetEventBlogById(ctx context.Context, id primitive.ObjectID) (*domain.EventBlog, error)
	GetAllEventBlogs(ctx context.Context) ([]domain.EventBlog, error)
	UpdateEventBlog(ctx context.Context, eventBlog domain.EventBlog) (*domain.EventBlog, error)
	DeleteEventBlog(ctx context.Context, id primitive.ObjectID) error
}

type eventBlogRepository struct {
	db              *mongo.Database
	collection_name string
}

func NewEventBlogRepository(db *mongo.Database, collection_name string) EventBlogRepository {
	return &eventBlogRepository{
		db:              db,
		collection_name: collection_name,
	}
}

// CreateEventBlog creates a new event blog entry
func (r *eventBlogRepository) CreateEventBlog(ctx context.Context, eventBlog domain.EventBlog) (*domain.EventBlog, error) {
	collection := r.db.Collection(r.collection_name)
	result, err := collection.InsertOne(ctx, eventBlog)
	if err != nil {
		return nil, err
	}

	// Setting the ID after insertion
	eventBlog.Id = result.InsertedID.(primitive.ObjectID)
	return &eventBlog, nil
}

// GetEventBlogById retrieves an event blog by its ID
func (r *eventBlogRepository) GetEventBlogById(ctx context.Context, id primitive.ObjectID) (*domain.EventBlog, error) {
	collection := r.db.Collection(r.collection_name)
	var eventBlog domain.EventBlog
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&eventBlog)
	if err != nil {
		return nil, err
	}
	return &eventBlog, nil
}

// GetAllEventBlogs retrieves all event blogs
func (r *eventBlogRepository) GetAllEventBlogs(ctx context.Context) ([]domain.EventBlog, error) {
	collection := r.db.Collection(r.collection_name)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var eventBlogs []domain.EventBlog
	for cursor.Next(ctx) {
		var eventBlog domain.EventBlog
		if err := cursor.Decode(&eventBlog); err != nil {
			return nil, err
		}
		eventBlogs = append(eventBlogs, eventBlog)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return eventBlogs, nil
}

// UpdateEventBlog updates an existing event blog
func (r *eventBlogRepository) UpdateEventBlog(ctx context.Context, eventBlog domain.EventBlog) (*domain.EventBlog, error) {
	collection := r.db.Collection(r.collection_name)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": eventBlog.Id}, bson.M{
		"$set": bson.M{
			"title":             eventBlog.Title,
			"short_description": eventBlog.ShortDescription,
			"created_at":        eventBlog.CreatedAt,
		},
	})
	if err != nil {
		return nil, err
	}
	return &eventBlog, nil
}

// DeleteEventBlog deletes an event blog by its ID
func (r *eventBlogRepository) DeleteEventBlog(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection(r.collection_name)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
