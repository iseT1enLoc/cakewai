package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	GoogleId       *string            `json:"google_id" bson:"google_id"`
	ProfilePicture *string            `json:"profile_picture" bson:"profile_picture"`
	Name           string             `json:"name" bson:"name"`
	Password       string             `json:"password" bson:"password"`
	Email          string             `json:"email" bson:"email"`
	Phone          string             `json:"phone" bson:"phone"`
	Address        Address            `json:"address" bson:"address"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      *time.Time         `json:"updated_at" bson:"updated_at"`
	IsAdmin        bool               `json:"is_admin,omitempty" bson:"is_admin,omitempty"  default:"false"`
	RoleID         string             `json:"role_id,omitempty" bson:"role_id,omitempty"`
}

//	type UserResponse struct {
//		Id             primitive.ObjectID `json:"id" bson:"_id"`
//		GoogleId       *string            `json:"google_id" bson:"google_id"`
//		ProfilePicture *string            `json:"profile_picture" bson:"profile_picture"`
//		Name           string             `json:"name" bson:"name"`
//		Email          string             `json:"email" bson:"email"`
//		Phone          string             `json:"phone" bson:"phone"`
//		Address        Address            `json:"address" bson:"address"`
//		IsAdmin        bool               `json:"is_admin" bson:"is_admin"`
//		CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
//		RoleID         string             `json:"role_id,omitempty" bson:"role_id,omitempty"`
//	}
type UserResponse struct {
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	GoogleId       *string            `json:"google_id" bson:"google_id"`
	ProfilePicture *string            `json:"profile_picture" bson:"profile_picture"`
	Name           string             `json:"name" bson:"name"`
	Email          string             `json:"email" bson:"email"`
	Phone          string             `json:"phone" bson:"phone"`
	Address        Address            `json:"address" bson:"address"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
}

type UserUseCase interface {
	GetUserById(c context.Context, id string) (*UserResponse, error)
	GetListUsers(c context.Context) ([]*UserResponse, error)
	UpdateUser(c context.Context, user *User) error
	DeleteUser(c context.Context, id string) error
}
