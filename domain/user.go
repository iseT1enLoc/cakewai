package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id             primitive.ObjectID `json:"id" bson:"id"`
	GoogleId       string             `json:"google_id" bson:"google_id"`
	ProfilePicture string             `json:"profile_picture" bson:"profile_picture"`
	Name           string             `json:"name" bson:"name"`
	Password       string             `json:"password" bsno:"password"`
	Email          string             `json:"email" bson:"email"`
	Phone          string             `json:"phone" bson:"phone"`
	Address        Address            `json:"address" bson:"address"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
	Invoices       []Invoice          `json:"invoices" bson:"invoices"`
	UserCart       []CartItem         `json:"cart_item" bson:"cart_item"`
	UserType       string             `json:"user_type" bson:"user_type"`
}

type UserResponse struct {
	Id             primitive.ObjectID `json:"id" bson:"id"`
	GoogleId       string             `json:"google_id" bson:"google_id"`
	ProfilePicture string             `json:"profile_picture" bson:"profile_picture"`
	Name           string             `json:"name" bson:"name"`
	Email          string             `json:"email" bson:"email"`
	Phone          string             `json:"phone" bson:"phone"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
}

type UserUseCase interface {
	GetUserById(c context.Context, id string) (*UserResponse, error)
	GetListUsers(c context.Context) ([]*UserResponse, error)
	UpdateUser(c context.Context, user *User) error
	DeleteUser(c context.Context, id string) error
}
