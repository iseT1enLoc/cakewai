package usecase

import (
	"context"
	"time"

	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/domain"
	tokenutil "cakewai/cakewai.com/internals/token_utils"
	"cakewai/cakewai.com/repository"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"
	"golang.org/x/crypto/bcrypt"
)

type signupUseCase struct {
	userRepository repository.UserRepository
	contextTimeout time.Duration
}

// SignUp implements domain.SignupUseCase.
func (s *signupUseCase) SignUp(ctx context.Context, request domain.SignupRequest, env *appconfig.Env) (accessToken string, refreshToken string, err error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Error(err)
		return
	}

	request.Password = string(encryptedPassword)

	user := &domain.User{
		Name:      request.Name,
		Password:  request.Password,
		Email:     request.Email,
		CreatedAt: time.Now(),
	}

	user, err = s.userRepository.CreateUser(ctx, user)
	if err != nil {
		log.Error(err)
		return
	}

	accessToken, err = tokenutil.CreateAccessToken(user, env.ACCESS_SECRET, env.ACCESS_TOK_EXP)
	if err != nil {
		log.Error(err)
		return
	}

	refreshToken, err = tokenutil.CreateRefreshToken(user, env.REFRESH_SECRET, env.REFRESH_TOK_EXP)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func NewSignupUseCase(userRepository repository.UserRepository, timeout time.Duration) domain.SignupUseCase {
	return &signupUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}
