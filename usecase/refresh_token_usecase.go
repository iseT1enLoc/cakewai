package usecase

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/domain"
	tokenutil "cakewai/cakewai.com/internals/token_utils"
	"cakewai/cakewai.com/repository"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"
)

type refreshTokenUsecase struct {
	refreshTokenRepository repository.RefreshTokenRepository
	contextTimeOut         time.Duration
}

// DeleteRefreshtoken implements domain.RefreshTokenUseCase.
func (r *refreshTokenUsecase) DeleteRefreshtoken(ctx context.Context, current_RT string, env *appconfig.Env) error {
	// Validate the input
	if current_RT == "" {
		fmt.Println("Refresh token is empty")
		return fmt.Errorf("refresh token cannot be empty")
	}

	// Call the repository layer to delete the refresh token
	err := r.refreshTokenRepository.DeleteRefreshtoken(ctx, current_RT, env)
	if err != nil {
		fmt.Printf("Error deleting refresh token: %v", err)
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	// Log success
	fmt.Printf("Successfully deleted refresh token: %s", current_RT)
	return nil
}

// RenewAccessToken implements domain.RefreshTokenUseCase.
func (r *refreshTokenUsecase) RenewAccessToken(ctx context.Context, refresh domain.RefreshTokenRequest, env *appconfig.Env) (access_token string, refresh_token string, err error) {
	re_token, err := r.GetRefreshTokenFromDB(ctx, refresh.RefreshToken, env)
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	uid, _ := primitive.ObjectIDFromHex(re_token.UserID)
	fmt.Printf("\nUserID of this app is that : %s\n", re_token.ID)
	newacc, _, _ := tokenutil.CreateAccessToken(uid, env.ACCESS_SECRET, false, 1000)
	return newacc, re_token.RefreshToken, nil

}

// GetRefreshTokenFromDB implements domain.RefreshTokenUseCase.
func (r *refreshTokenUsecase) GetRefreshTokenFromDB(ctx context.Context, current_refresh_token string, env *appconfig.Env) (*domain.RefreshTokenRequest, error) {
	token, err := r.refreshTokenRepository.GetRefreshTokenFromDB(ctx, current_refresh_token, env)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return token, nil
}

// InsertRefreshTokenToDB implements domain.RefreshTokenUseCase.
func (r *refreshTokenUsecase) InsertRefreshTokenToDB(ctx context.Context, refresh_token domain.RefreshTokenRequest, user_id string, is_admin bool, env *appconfig.Env) (string, error) {
	res, err := r.refreshTokenRepository.InsertRefreshTokenToDB(ctx, user_id, is_admin, env)
	print("enter refresh usecase")
	if err != nil {
		log.Error(err)
		return res, err
	}
	return res, nil
}

// RevokeToken implements domain.RefreshTokenUseCase.
func (r *refreshTokenUsecase) RevokeToken(ctx context.Context, current_RT string, env *appconfig.Env) error {
	err := r.refreshTokenRepository.RevokeToken(ctx, current_RT, env)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// UpdateRefreshTokenChanges implements domain.RefreshTokenUseCase.
func (r *refreshTokenUsecase) UpdateRefreshTokenChanges(ctx context.Context, updatedRT domain.RefreshTokenRequest, env *appconfig.Env) (*domain.RefreshTokenRequest, error) {
	res, err := r.refreshTokenRepository.UpdateRefreshTokenChanges(ctx, updatedRT, env)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return res, nil
}

// RefreshToken implements domain.RefreshTokenUseCase.
func (r *refreshTokenUsecase) RefreshToken(ctx context.Context, request domain.RefreshTokenRequest, is_admin bool, env *appconfig.Env) (accessToken string, refreshToken string, err error) {
	//FOCUS------------
	_, _, errs := tokenutil.ExtractIDAndRole(request.RefreshToken, env.REFRESH_SECRET)
	if errs != nil {
		fmt.Println("line 25 Oh my godness")
		log.Error(err)

		return
	}

	// return accessToken, refreshToken, nil
	accesstoken, refresh_token, err := r.refreshTokenRepository.RefreshToken(ctx, request.RefreshToken, is_admin, env)
	if err != nil {
		log.Error(err)
		return "", "", err
	}
	fmt.Println("line 101 happy happy")

	return accesstoken, refresh_token, nil
}

func NewRefreshTokenUseCase(refreshTokenRepository repository.RefreshTokenRepository, timeout time.Duration) domain.RefreshTokenUseCase {
	return &refreshTokenUsecase{
		refreshTokenRepository: refreshTokenRepository,
		contextTimeOut:         timeout,
	}
}
