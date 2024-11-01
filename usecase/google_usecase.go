package usecase

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	apperror "cakewai/cakewai.com/component/apperr"
	"cakewai/cakewai.com/domain"
	tokenutil "cakewai/cakewai.com/internals/token_utils"
	"cakewai/cakewai.com/repository"

	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"

	"golang.org/x/oauth2"
)

type googleUseCase struct {
	userRepository repository.UserRepository
	contextTimeout time.Duration
}

func NewGoogleUseCase(userRepository repository.UserRepository, timeout time.Duration) domain.GoogleUseCase {
	return &googleUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (lu *googleUseCase) GoogleLogin(ctx context.Context, data []byte, env *appconfig.Env) (accessToken string, refreshToken string, err error) {
	var googleUser *domain.GoogleUser
	err = json.Unmarshal(data, &googleUser)
	if err != nil {
		log.Error(err)
		return
	}

	user := &domain.User{
		GoogleId:       googleUser.Id,
		ProfilePicture: googleUser.Picture,
		Email:          googleUser.Email,
		Name:           googleUser.Name,
	}

	var existingUser *domain.User
	existingUser, err = lu.userRepository.GetUserByEmail(ctx, googleUser.Email)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			user, err = lu.userRepository.CreateUser(ctx, user)
			if err != nil {
				log.Error(err)
				return
			}
		} else {
			log.Error(err)
			return
		}
	}

	if existingUser != nil {
		user = existingUser
	}

	// Create access token
	accessToken, err = tokenutil.CreateAccessToken(user.Id, env.ACCESS_SECRET, int(env.ACCESS_TOK_EXP))
	if err != nil {
		log.Error(err)
		return
	}

	// Create refresh token
	refreshToken, err = tokenutil.CreateRefreshToken(user.Id, env.REFRESH_SECRET, env.REFRESH_TOK_EXP)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func (lu *googleUseCase) GetUserDataFromGoogle(googleOauthConfig *oauth2.Config, code, oauthGoogleUrlAPI string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Error(err)
		return nil, apperror.ErrCodeExchangeWrong
	}

	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		log.Error(err)
		return nil, apperror.ErrFailedGetGoogleUser
	}

	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return nil, apperror.ErrFailedToReadResponse
	}

	return contents, nil
}

func (lu *googleUseCase) GenerateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}
