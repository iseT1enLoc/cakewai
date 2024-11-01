package repository

import (
	"context"
	"time"

	"cakewai/cakewai.com/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RefreshTokenRepository interface {
	RefreshToken(ctx context.Context, current_RT string) (string, error)
	RevokeToken(ctx context.Context, current_RT string) error
	InsertRefreshTokenToDB(ctx context.Context, refresh_token domain.RefreshTokenRequest, user_id string) error
	GetRefreshTokenFromDB(ctx context.Context, current_refresh_token string) (*domain.RefreshTokenRequest, error)
	UpdateRefreshTokenChanges(ctx context.Context, updatedRT domain.RefreshTokenRequest) error
}

type refreshtokenRepository struct {
	db              *mongo.Database
	collection_name string
}

// UpdateRefreshTokenChanges implements RefreshTokenRepository.
func (r *refreshtokenRepository) UpdateRefreshTokenChanges(ctx context.Context, updatedRT domain.RefreshTokenRequest) error {
	collection := r.db.Collection(r.collection_name)
	_, err := collection.UpdateByID(ctx, bson.M{"_id": updatedRT.ID}, updatedRT)
	if err != nil {
		return err
	}
	return nil
}

// GetRefreshTokenFromDB implements RefreshTokenRepository.
func (r *refreshtokenRepository) GetRefreshTokenFromDB(ctx context.Context, current_refresh_token string) (*domain.RefreshTokenRequest, error) {
	collection := r.db.Collection(r.collection_name)

	var refreshToken domain.RefreshTokenRequest
	err := collection.FindOne(ctx, bson.M{"refresh_token": current_refresh_token}).Decode(&refreshToken)
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// RefreshToken implements RefreshTokenRepository.
func (r *refreshtokenRepository) RefreshToken(ctx context.Context, current_RT string) (string, error) {
	//get current refresh token from database
	refresh_token, _ := r.GetRefreshTokenFromDB(ctx, current_RT)

	//create old refresh token and save the new one
	refresh_token.RevokeAt = time.Now().Local()
	refresh_token.ReplaceByRT = current_RT
}

// RevokeToken implements RefreshTokenRepository.
func (r *refreshtokenRepository) RevokeToken(ctx context.Context, current_RT string) error {
	RT, err := r.GetRefreshTokenFromDB(ctx, current_RT)

	if err != nil {
		return err
	}

	RT.RevokeAt = time.Now()
	RT.IsActive = false
	RT.IsExpire = true
	err = r.UpdateRefreshTokenChanges(ctx, *RT)
	if err != nil {
		return err
	}
	return nil
}

// SaveRefreshTokenToDB implements RefreshTokenRepository.
func (r *refreshtokenRepository) InsertRefreshTokenToDB(ctx context.Context, refresh_token domain.RefreshTokenRequest, user_id string) error {
	collection := r.db.Collection(r.collection_name)
	_, err := collection.InsertOne(ctx, bson.M{"_id": primitive.NewObjectID(),
		"RefreshToken": refresh_token.RefreshToken,
		"UserID":       refresh_token.UserID,
		"ExpireAt":     refresh_token.ExpireAt,
		"CreatedAt":    refresh_token.CreatedAt,
		"RevokeAt":     refresh_token.RevokeAt,
		"ReplaceByRT":  refresh_token.ReplaceByRT,
		"IsActive":     refresh_token.IsActive,
		"IsExpire":     refresh_token.IsActive})
	if err != nil {
		return err
	}
	return nil
}

func NewrefreshTokenRepository(db *mongo.Database, collection_name string) RefreshTokenRepository {
	return &refreshtokenRepository{
		db:              db,
		collection_name: collection_name,
	}
}
