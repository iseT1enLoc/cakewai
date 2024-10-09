package usecase

import (
	"cakewai/cakewai.com/domain"
	"cakewai/cakewai.com/repository"
	"context"
	"time"
)

type profileUseCase struct {
	userRepository repository.UserRepository
	contextTimeOut time.Duration
}

func NewProfileUseCase(time_expiration_time time.Duration, user_repository repository.UserRepository) domain.UserUseCase {
	return &profileUseCase{
		userRepository: user_repository,
		contextTimeOut: 0,
	}

}

// DeleteUser implements domain.UserUseCase.
func (p *profileUseCase) DeleteUser(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, p.contextTimeOut)
	defer cancel()
	return p.userRepository.DeleteUser(ctx, id)
}

// GetListUsers implements domain.UserUseCase.
func (p *profileUseCase) GetListUsers(c context.Context) ([]*domain.UserResponse, error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeOut)
	defer cancel()
	var urs []*domain.UserResponse
	users, err := p.userRepository.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		urs = append(urs, &domain.UserResponse{
			Id:             user.Id,
			GoogleId:       user.GoogleId,
			ProfilePicture: user.ProfilePicture,
			Name:           user.Name,
			Email:          user.Email,
			Phone:          user.Phone,
			CreatedAt:      user.CreatedAt,
		})
	}
	return urs, nil
}

// GetUserById implements domain.UserUseCase.
func (p *profileUseCase) GetUserById(c context.Context, id string) (*domain.UserResponse, error) {
	ctx, cancel := context.WithTimeout(c, time.Hour*5)
	defer cancel()
	var ur *domain.UserResponse
	user, err := p.userRepository.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	ur = &domain.UserResponse{
		Id:             user.Id,
		GoogleId:       user.GoogleId,
		ProfilePicture: user.ProfilePicture,
		Name:           user.Name,
		Email:          user.Email,
		Phone:          user.Phone,
		CreatedAt:      user.CreatedAt,
	}
	return ur, nil
}

// UpdateUser implements domain.UserUseCase.
func (p *profileUseCase) UpdateUser(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, p.contextTimeOut)
	defer cancel()
	return p.userRepository.UpdateUser(ctx, user)
}
