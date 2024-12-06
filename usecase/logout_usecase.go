package usecase

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/domain"
	"cakewai/cakewai.com/repository"
	"context"
	"fmt"
	"time"
)

type logoutUsecase struct {
	refreshRepository repository.RefreshTokenRepository
	contextTimeout    time.Duration
}

// LogOut implements domain.LogOutUseCase.
func (l *logoutUsecase) LogOut(ctx context.Context, request domain.LogoutRequest, env *appconfig.Env) (err error) {
	// Validate the input
	if request.Refresh_token == "" {
		fmt.Println("Refresh token is empty")
		return fmt.Errorf("refresh token cannot be empty")
	}

	// Set a timeout for the context
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	// Attempt to delete the refresh token
	err = l.refreshRepository.DeleteRefreshtoken(ctx, request.Refresh_token, env)
	if err != nil {
		fmt.Printf("Failed to delete refresh token: %v", err)
		return fmt.Errorf("failed to log out: %w", err)
	}

	// Log success
	fmt.Printf("User successfully logged out and refresh token deleted: %s", request.Refresh_token)
	return nil
}

func NewLogoutUseCase(refreshTokenRepository repository.RefreshTokenRepository, timeout time.Duration) domain.LogOutUseCase {
	return &logoutUsecase{
		refreshRepository: refreshTokenRepository,
		contextTimeout:    timeout,
	}
}
