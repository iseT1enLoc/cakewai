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
	userRepository repository.UserRepository
	contextTimeOut time.Duration
}

// Login implements domain.LoginUseCase.
func (l *loginUsecase) Login(ctx context.Context, request domain.LoginRequest, env *appconfig.Env) (accessToken string, refreshToken string, err error) {
	var user *domain.User
	fmt.Print("line 25 login usecase")
	user, err = l.userRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		log.Error(err)
		return
	}

	if user.GoogleId != "" {
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
	arg := "hihi"
	accessToken, err = tokenutil.CreateAccessToken(user, arg, time.Now().Second()*3600)
	fmt.Printf("\n Create Access token %s\n", arg)
	if err != nil {
		log.Error(err)
		return
	}

	refreshToken, err = tokenutil.CreateRefreshToken(user, env.REFRESH_SECRET, env.REFRESH_TOK_EXP)
	if err != nil {
		log.Error(err)
		return
	}

	return accessToken, refreshToken, nil
}

func NewLoginUseCase(user repository.UserRepository, timeout time.Duration) domain.LoginUseCase {
	return &loginUsecase{
		userRepository: user,
		contextTimeOut: timeout,
	}
}
