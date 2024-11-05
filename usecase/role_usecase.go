package usecase

import (
	"cakewai/cakewai.com/domain"
	"cakewai/cakewai.com/repository"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type roleUsecase struct {
	roleRepository repository.RoleRepository
	contextTimeout time.Duration
}

// GetRoleByRoleName implements domain.RoleUsecase.
func (r *roleUsecase) GetRoleByRoleName(ctx context.Context, rolename string) (*domain.Role, error) {
	role, err := r.roleRepository.GetRoleByRoleName(ctx, rolename)
	return role, err
}

// CreateRole implements domain.RoleUsecase.
func (r *roleUsecase) CreateRole(ctx context.Context, role domain.Role) error {
	err := r.roleRepository.CreateRole(ctx, role)
	return err
}

// GetAllRoles implements domain.RoleUsecase.
func (r *roleUsecase) GetAllRoles(ctx context.Context) (roles []*domain.Role, err error) {
	roles, err = r.roleRepository.GetAllRoles(ctx)
	return roles, err
}

// GetRoleByID implements domain.RoleUsecase.
func (r *roleUsecase) GetRoleByID(ctx context.Context, id primitive.ObjectID) (role *domain.Role, err error) {
	role, err = r.roleRepository.GetRoleByID(ctx, id)
	if err != nil {
		return role, err
	}
	return role, nil
}

// UpdateRole implements domain.RoleUsecase.
func (r *roleUsecase) UpdateRole(ctx context.Context, updatedRole domain.Role) (*domain.Role, error) {
	_, err := r.roleRepository.UpdateRole(ctx, updatedRole)
	if err != nil {
		return nil, err
	}
	return &updatedRole, nil
}

func NewRoleUsecase(rolerepo repository.RoleRepository, time time.Duration) domain.RoleUsecase {
	return &roleUsecase{
		roleRepository: rolerepo,
		contextTimeout: time,
	}
}
