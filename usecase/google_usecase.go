package usecase

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	apperror "cakewai/cakewai.com/component/apperr"
	"cakewai/cakewai.com/domain"
	tokenutil "cakewai/cakewai.com/internals/token_utils"
	"cakewai/cakewai.com/repository"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"context"
	"crypto/rand"
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
	refresh_repo   repository.RefreshTokenRepository
	contextTimeout time.Duration
}

func NewGoogleUseCase(userRepository repository.UserRepository, refreshTokenRepository repository.RefreshTokenRepository, timeout time.Duration) domain.GoogleUseCase {
	return &googleUseCase{
		userRepository: userRepository,
		refresh_repo:   refreshTokenRepository,
		contextTimeout: timeout,
	}
}

func (lu *googleUseCase) GoogleLogin(ctx context.Context, data []byte, env *appconfig.Env) (accessToken string, refreshToken string, err error) {
	var googleUser *domain.GoogleUser
	fmt.Println("Enter line 39 google login")
	err = json.Unmarshal(data, &googleUser)
	fmt.Println("Enter line 41 google login")
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println("Enter line 46 google login")
	user := &domain.User{
		Id:             primitive.NewObjectID(),
		GoogleId:       &googleUser.Id,
		ProfilePicture: &googleUser.Picture,
		Email:          googleUser.Email,
		Name:           googleUser.Name,
		CreatedAt:      time.Now().UTC(),
	}

	fmt.Println("Enter line 53 google login")

	var existingUser *domain.User
	existingUser, err = lu.userRepository.GetUserByEmail(ctx, googleUser.Email)
	if err != nil {
		fmt.Println(err)

		user, err = lu.userRepository.CreateUser(ctx, user)

		fmt.Println("Enter line 61 google login")
		fmt.Println(err)
		if err != nil {
			log.Error(err)
			return
		}
	}
	fmt.Println("Enter line 69 google login")
	if existingUser != nil {
		user = existingUser
	}

	accessToken, _, err = tokenutil.CreateAccessToken(user.Id, env.ACCESS_SECRET, false, env.ACCESS_TOK_EXP)
	if err != nil {
		log.Error(err)
		return
	}
	uidstring := user.Id.Hex()

	refreshToken, err = lu.refresh_repo.InsertRefreshTokenToDB(ctx, uidstring, user.IsAdmin, env)
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
