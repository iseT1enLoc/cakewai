package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	RoleID     primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `json:"role_name" bson:"role_name"`
	Decription string             `json:"description" bson:"description"`
}

type RoleUsecase interface {
	CreateRole(ctx context.Context, role Role) error
	GetAllRoles(ctx context.Context) (roles []*Role, err error)
	GetRoleByID(ctx context.Context, id primitive.ObjectID) (role *Role, err error)
	UpdateRole(ctx context.Context, updatedRole Role) (*Role, error)
	GetRoleByRoleName(ctx context.Context, rolename string) (*Role, error)
}
