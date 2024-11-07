package usecase

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/domain"
	tokenutil "cakewai/cakewai.com/internals/token_utils"
	"cakewai/cakewai.com/repository"
	"context"
	"fmt"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"
)

type refreshTokenUsecase struct {
	refreshTokenRepository repository.RefreshTokenRepository
	contextTimeOut         time.Duration
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
func (r *refreshTokenUsecase) InsertRefreshTokenToDB(ctx context.Context, refresh_token domain.RefreshTokenRequest, user_id string, env *appconfig.Env) (string, error) {
	res, err := r.refreshTokenRepository.InsertRefreshTokenToDB(ctx, user_id, env)
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
func (r *refreshTokenUsecase) RefreshToken(ctx context.Context, request domain.RefreshTokenRequest, currentRT string, env *appconfig.Env) (accessToken string, refreshToken string, err error) {
	//FOCUS------------
	_, errs := tokenutil.ExtractID(currentRT, env.REFRESH_SECRET)
	if errs != nil {
		fmt.Println("line 25 Oh my godness")
		log.Error(err)

		return
	}
	// fmt.Println(id)
	// var user *domain.User
	// //id = "672167f95b55a7c71f00f18a"
	// fmt.Println("Enter line 30 refresh token usecase line 30")
	// user, err = r.userRepository.GetUserById(ctx, id)

	// fmt.Println("Enter line 31 refresh token usecase line 31")
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }
	// fmt.Println("Enter refresh token usecase line 35")
	// accessToken, err = tokenutil.CreateAccessToken(user.Id, env.ACCESS_SECRET, env.ACCESS_TOK_EXP)
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }
	// fmt.Println("Enter refresh token usecase linen 41")
	// refreshToken, err = tokenutil.CreateRefreshToken(user.Id, env.REFRESH_SECRET, env.REFRESH_TOK_EXP)
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }

	// return accessToken, refreshToken, nil
	accesstoken, err := r.refreshTokenRepository.RefreshToken(ctx, currentRT, env)
	if err != nil {
		log.Error(err)
		return "", "", err
	}
	fmt.Println("line 101 happy happy")
	// _, err = r.UpdateRefreshTokenChanges(ctx, *RT, env)
	// if err != nil {
	// 	log.Error(err)
	// 	return "", "", err
	// }
	return accesstoken, currentRT, nil
}

func NewRefreshTokenUseCase(refreshTokenRepository repository.RefreshTokenRepository, timeout time.Duration) domain.RefreshTokenUseCase {
	return &refreshTokenUsecase{
		refreshTokenRepository: refreshTokenRepository,
		contextTimeOut:         timeout,
	}
}
