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
	userRepository repository.UserRepository
	contextTimeOut time.Duration
}

// RefreshToken implements domain.RefreshTokenUseCase.
func (r *refreshTokenUsecase) RefreshToken(ctx context.Context, request domain.RefreshTokenRequest, env *appconfig.Env) (accessToken string, refreshToken string, err error) {
	var id string
	id, err = tokenutil.ExtractIDFromToken(request.RefreshToken, env.REFRESH_SECRET)
	if err != nil {
		fmt.Println("line 25 Oh my godness")
		log.Error(err)
		return
	}

	var user *domain.User
	fmt.Println("Enter line 30 refresh token usecase line 30")
	user, err = r.userRepository.GetUserById(ctx, id)

	fmt.Println("Enter line 31 refresh token usecase line 31")
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println("Enter refresh token usecase line 35")
	accessToken, err = tokenutil.CreateAccessToken(user, env.ACCESS_SECRET, env.ACCESS_TOK_EXP)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println("Enter refresh token usecase linen 41")
	refreshToken, err = tokenutil.CreateRefreshToken(user, env.REFRESH_SECRET, env.REFRESH_TOK_EXP)
	if err != nil {
		log.Error(err)
		return
	}

	return accessToken, refreshToken, nil
}

func NewRefreshTokenUseCase(userrepository repository.UserRepository, timeout time.Duration) domain.RefreshTokenUseCase {
	return &refreshTokenUsecase{
		userRepository: userrepository,
		contextTimeOut: timeout,
	}
}
