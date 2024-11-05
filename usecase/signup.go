package usecase

import (
	"context"
	"time"

	appconfig "cakewai/cakewai.com/component/appcfg"
	apperror "cakewai/cakewai.com/component/apperr"
	"cakewai/cakewai.com/domain"
	tokenutil "cakewai/cakewai.com/internals/token_utils"
	"cakewai/cakewai.com/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"
	"golang.org/x/crypto/bcrypt"
)

type signupUseCase struct {
	userRepository   repository.UserRepository
	refreshTokenrepo repository.RefreshTokenRepository
	contextTimeout   time.Duration
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
	//check if email have already sign up
	if _, err := s.userRepository.GetUserByEmail(ctx, request.Email); err == nil {
		accessToken = ""
		refreshToken = ""
		var ErrEmailAlreadyRegistered = apperror.ErrEmailAlreadyExist
		log.Error(ErrEmailAlreadyRegistered)
		return "", "", ErrEmailAlreadyRegistered
	}

	request.Password = string(encryptedPassword)

	user := &domain.User{
		Id:        primitive.NewObjectID(),
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

	accessToken, err = tokenutil.CreateAccessToken(user.Id, env.ACCESS_SECRET, env.ACCESS_TOK_EXP)
	if err != nil {
		log.Error(err)
		return
	}
	uidstring := user.Id.Hex()
	print(uidstring)
	refreshToken, err = s.refreshTokenrepo.InsertRefreshTokenToDB(ctx, uidstring, env)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func NewSignupUseCase(userRepository repository.UserRepository, refreshtokenRepo repository.RefreshTokenRepository, timeout time.Duration) domain.SignupUseCase {
	return &signupUseCase{
		userRepository:   userRepository,
		refreshTokenrepo: refreshtokenRepo,
		contextTimeout:   timeout,
	}
}
