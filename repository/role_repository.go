package repository

import (
	"context"
	"strings"

	"cakewai/cakewai.com/domain"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type roleRepository struct {
	db              *mongo.Database
	collection_name string
}

// GetRoleByRoleName implements RoleRepository.
func (r *roleRepository) GetRoleByRoleName(ctx context.Context, role_name string) (*domain.Role, error) {
	role_name = strings.ToLower(role_name)
	collection := r.db.Collection(r.collection_name)
	var role *domain.Role
	err := collection.FindOne(ctx, bson.M{"role_name": role_name}).Decode(&role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// CreateRole implements RoleRepository.
func (r *roleRepository) CreateRole(ctx context.Context, role domain.Role) error {
	collection := r.db.Collection(r.collection_name)
	_, err := collection.InsertOne(ctx, bson.M{"_id": role.RoleID, "role_name": role.Name, "description": role.Decription})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil

}

// GetAllRoles implements RoleRepository.
func (r *roleRepository) GetAllRoles(ctx context.Context) (roles []*domain.Role, err error) {
	collection := r.db.Collection(r.collection_name)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleByID implements RoleRepository.
func (r *roleRepository) GetRoleByID(ctx context.Context, id primitive.ObjectID) (role *domain.Role, err error) {
	collection := r.db.Collection(r.collection_name)
	var expectedRole domain.Role
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&expectedRole)
	if err != nil {
		return nil, err
	}
	return &expectedRole, err
}

// UpdateRole implements RoleRepository.
func (r *roleRepository) UpdateRole(ctx context.Context, role domain.Role) (*domain.Role, error) {
	collection := r.db.Collection(r.collection_name)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": role.RoleID},
		bson.M{"$set": bson.M{"role_name": role.Name, "description": role.Decription}})
	return &role, err
}

type RoleRepository interface {
	CreateRole(ctx context.Context, role domain.Role) error
	GetAllRoles(ctx context.Context) (roles []*domain.Role, err error)
	GetRoleByID(ctx context.Context, id primitive.ObjectID) (role *domain.Role, err error)
	UpdateRole(ctx context.Context, role domain.Role) (*domain.Role, error)
	GetRoleByRoleName(ctx context.Context, role_name string) (*domain.Role, error)
}

func NewRoleRepository(db *mongo.Database, collection_name string) RoleRepository {
	return &roleRepository{
		db:              db,
		collection_name: collection_name,
	}
}
