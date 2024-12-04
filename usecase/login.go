package usecase

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	apperror "cakewai/cakewai.com/component/apperr"
	"cakewai/cakewai.com/domain"
	tokenutil "cakewai/cakewai.com/internals/token_utils"
	"cakewai/cakewai.com/repository"
	"context"
	"fmt"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"
	"golang.org/x/crypto/bcrypt"
)

type loginUsecase struct {
	userRepository   repository.UserRepository
	refreshTokenrepo repository.RefreshTokenRepository
	contextTimeOut   time.Duration
}

// Login implements domain.LoginUseCase.
func (l *loginUsecase) Login(ctx context.Context, request domain.LoginRequest, env *appconfig.Env) (accessToken string, refreshToken string, err error) {
	var user *domain.User
	fmt.Print("line 25 login usecase")
	user, err = l.userRepository.GetUserByEmail(ctx, request.Email)
	fmt.Printf("\nUser id of this function is that %v", user.Id)
	if err != nil {
		log.Error(err)
		return
	}

	if user.GoogleId != nil {
		log.Error(err)
		err = apperror.ErrUserShouldLoginWithGoogle
		return
	}
	fmt.Print("line 36 login usecase")

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		fmt.Printf("Password retrieve from database %s\n", user.Password)
		fmt.Println(request.Password)
		log.Error(err)
		err = apperror.ErrInvalidPassword
		return
	}
	fmt.Print("line 43 login usecase")
	accessToken, err = tokenutil.CreateAccessToken(user.Id, env.ACCESS_SECRET, time.Now().Second()*3600)
	fmt.Printf("\n Create Access token %s\n", env.ACCESS_SECRET)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println("Token id from line 55")
	//refreshToken, err = tokenutil.CreateRefreshToken(user.Id, env.REFRESH_SECRET, env.REFRESH_TOK_EXP)

	refreshtoken, err := l.refreshTokenrepo.InsertRefreshTokenToDB(ctx, user.Id.Hex(), env)
	fmt.Println("line 59 vui ve vui ve")
	if err != nil {
		log.Error(err)
		return
	}

	return accessToken, refreshtoken, nil
}

func NewLoginUseCase(user repository.UserRepository, refreshtoken repository.RefreshTokenRepository, timeout time.Duration) domain.LoginUseCase {
	return &loginUsecase{
		userRepository:   user,
		refreshTokenrepo: refreshtoken,
		contextTimeOut:   timeout,
	}
}
