package repository

import (
	"context"
	"fmt"
	"time"

	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/domain"
	tokenutil "cakewai/cakewai.com/internals/token_utils"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RefreshTokenRepository interface {
	RefreshToken(ctx context.Context, current_RT string, env *appconfig.Env) (accesstoken string, err error)
	RevokeToken(ctx context.Context, current_RT string, env *appconfig.Env) error
	InsertRefreshTokenToDB(ctx context.Context, user_id string, env *appconfig.Env) (string, error)
	GetRefreshTokenFromDB(ctx context.Context, current_refresh_token string, env *appconfig.Env) (*domain.RefreshTokenRequest, error)
	UpdateRefreshTokenChanges(ctx context.Context, updatedRT domain.RefreshTokenRequest, env *appconfig.Env) (*domain.RefreshTokenRequest, error)
}

type refreshtokenRepository struct {
	db              *mongo.Database
	collection_name string
}

// UpdateRefreshTokenChanges implements RefreshTokenRepository.
func (r *refreshtokenRepository) UpdateRefreshTokenChanges(ctx context.Context, updatedRT domain.RefreshTokenRequest, env *appconfig.Env) (*domain.RefreshTokenRequest, error) {
	collection := r.db.Collection(r.collection_name)

	oid := primitive.ObjectIDFromHex(updateRT.ID)
	// Define the filter by _id
	filter := bson.M{"_id": updateRT.ID}

	// Use $set to specify which fields to update
	update := bson.M{
		"$set": newValues,
	}

	// Perform the update operation
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &updatedRT, nil
}

// GetRefreshTokenFromDB implements RefreshTokenRepository.
func (r *refreshtokenRepository) GetRefreshTokenFromDB(ctx context.Context, current_refresh_token string, env *appconfig.Env) (*domain.RefreshTokenRequest, error) {
	collection := r.db.Collection(r.collection_name)

	var refreshToken domain.RefreshTokenRequest
	err := collection.FindOne(ctx, bson.M{"refresh_token": current_refresh_token}).Decode(&refreshToken)
	fmt.Println("This is line 48 refresh token")
	if err != nil {
		// If no document is found, handle it appropriately
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("refresh token not found")
		}
		// For other errors, return as is
		return nil, err
	}

	fmt.Println(refreshToken)
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// RefreshToken implements RefreshTokenRepository.
func (r *refreshtokenRepository) RefreshToken(ctx context.Context, current_RT string, env *appconfig.Env) (accesstoken string, err error) {
	//get current refresh token from database
	refresh_token, err := r.GetRefreshTokenFromDB(ctx, current_RT, env)
	if err != nil {
		log.Error(err)
		return "", err
	}
	id, _ := primitive.ObjectIDFromHex(refresh_token.ID.Hex())
	token, err := tokenutil.CreateAccessToken(id, env.ACCESS_SECRET, int(time.Second)*3600)

	if err != nil {
		return "", err
	}
	return token, nil
}

// RevokeToken implements RefreshTokenRepository.
func (r *refreshtokenRepository) RevokeToken(ctx context.Context, current_RT string, env *appconfig.Env) error {
	RT, err := r.GetRefreshTokenFromDB(ctx, current_RT, env)

	if err != nil {
		return err
	}

	RT.RevokeAt = time.Now()
	RT.IsActive = false
	RT.IsExpire = true
	_, err = r.UpdateRefreshTokenChanges(ctx, *RT, env)
	if err != nil {
		return err
	}
	return nil
}

// SaveRefreshTokenToDB implements RefreshTokenRepository.
func (r *refreshtokenRepository) InsertRefreshTokenToDB(ctx context.Context, user_id string, env *appconfig.Env) (string, error) {
	fmt.Println("Insert Refresh token to db")
	collection := r.db.Collection(r.collection_name)
	fmt.Println("line 100 insert rt repository")
	oid, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		log.Error(err)
		return "", err
	}
	refresh_token, err := tokenutil.CreateRefreshToken(oid, env.REFRESH_SECRET, int(time.Second)*3600)
	fmt.Println(refresh_token)
	fmt.Println("line 107 insert rt repository")
	refreshtoken := domain.RefreshTokenRequest{
		ID:           primitive.NewObjectID(),
		RefreshToken: refresh_token,
		UserID:       user_id,
		ExpireAt:     time.Now().Local().Add(time.Second * 3600),
		CreatedAt:    time.Now().Local(),
		IsActive:     true,
		IsExpire:     false,
	}
	_, err = collection.InsertOne(ctx, refreshtoken)
	fmt.Println("insert refresh token")
	if err != nil {
		log.Error(err)
		return "", err
	}
	return refresh_token, nil
}

func NewrefreshTokenRepository(db *mongo.Database, collection_name string) RefreshTokenRepository {
	return &refreshtokenRepository{
		db:              db,
		collection_name: collection_name,
	}
}
