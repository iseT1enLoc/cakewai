package repository

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/domain"
	tokenutil "cakewai/cakewai.com/internals/token_utils"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RefreshTokenRepository interface {
	RefreshToken(ctx context.Context, current_RT string, is_admin bool, env *appconfig.Env) (accesstoken string, refresh_token string, err error)
	RevokeToken(ctx context.Context, current_RT string, env *appconfig.Env) error
	InsertRefreshTokenToDB(ctx context.Context, user_id string, is_admin bool, env *appconfig.Env) (string, error)
	GetRefreshTokenFromDB(ctx context.Context, current_refresh_token string, env *appconfig.Env) (*domain.RefreshTokenRequest, error)
	UpdateRefreshTokenChanges(ctx context.Context, updatedRT domain.RefreshTokenRequest, env *appconfig.Env) (*domain.RefreshTokenRequest, error)
	DeleteRefreshtoken(ctx context.Context, current_RT string, env *appconfig.Env) error
	// Method for cleanup of expired tokens
	CleanupExpiredTokens(ctx context.Context) error
}

type refreshtokenRepository struct {
	db              *mongo.Database
	collection_name string
}

// DeleteRefreshtoken implements RefreshTokenRepository.
func (r *refreshtokenRepository) DeleteRefreshtoken(ctx context.Context, current_RT string, env *appconfig.Env) error {
	// Define the collection
	collection := r.db.Collection(r.collection_name) // Replace with your actual collection name

	// Construct the filter to identify the document to delete
	filter := bson.M{"refresh_token": current_RT}

	// Perform the delete operation
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Printf("Failed to delete refresh token: %v", err)
		return fmt.Errorf("could not delete refresh token: %w", err)
	}

	// Check if no documents were deleted
	if result.DeletedCount == 0 {
		fmt.Printf("No refresh token found to delete for token: %s", current_RT)
		return fmt.Errorf("refresh token not found")
	}

	// Log success
	fmt.Printf("Successfully deleted refresh token: %s", current_RT)
	return nil
}

// CleanupExpiredTokens implements RefreshTokenRepository.
func (r *refreshtokenRepository) CleanupExpiredTokens(ctx context.Context) error {
	collection := r.db.Collection(r.collection_name)

	// Delete tokens that are expired
	_, err := collection.DeleteMany(ctx, bson.M{
		"expire_at": bson.M{"$lt": time.Now()},
	})
	if err != nil {
		return fmt.Errorf("failed to delete expired tokens: %v", err)
	}

	return nil
}

// UpdateRefreshTokenChanges implements RefreshTokenRepository.
func (r *refreshtokenRepository) UpdateRefreshTokenChanges(ctx context.Context, updatedRT domain.RefreshTokenRequest, env *appconfig.Env) (*domain.RefreshTokenRequest, error) {
	collection := r.db.Collection(r.collection_name)

	// Convert the refresh token ID to ObjectID
	oid, err := primitive.ObjectIDFromHex(updatedRT.ID.Hex())
	if err != nil {
		log.Printf("Error converting ID to ObjectID: %v", err)
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	// Define the filter by _id
	filter := bson.M{"_id": oid}

	// Use $set to specify the fields to update
	update := bson.M{
		"$set": bson.M{
			"revoke_at":      time.Now().Local(),
			"replaced_token": updatedRT.RefreshToken,
			"is_active":      false,
			"is_expire":      true,
		},
	}

	// Perform the update operation
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating refresh token changes for user %s: %v", updatedRT.UserID, err)
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}

	// Check if any document was updated
	if result.MatchedCount == 0 {
		log.Printf("No refresh token found with ID: %s", updatedRT.ID.Hex())
		return nil, fmt.Errorf("no refresh token found for the provided ID: %s", updatedRT.ID.Hex())
	}

	// Log success
	log.Printf("Successfully updated refresh token for user %s", updatedRT.UserID)

	// Return the updated refresh token
	return &updatedRT, nil
}

// GetRefreshTokenFromDB implements RefreshTokenRepository.
func (r *refreshtokenRepository) GetRefreshTokenFromDB(ctx context.Context, current_refresh_token string, env *appconfig.Env) (*domain.RefreshTokenRequest, error) {
	// Set a timeout for the database query
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	collection := r.db.Collection(r.collection_name)

	var refreshToken domain.RefreshTokenRequest
	err := collection.FindOne(ctx, bson.M{"refresh_token": current_refresh_token}).Decode(&refreshToken)

	// Enhanced logging for better traceability
	if err != nil {
		log.Printf("Error occurred while fetching refresh token: %v", err)

		// If no document is found, return a custom error
		if err == mongo.ErrNoDocuments {
			log.Printf("Refresh token not found for token: %s", current_refresh_token)
			return nil, fmt.Errorf("refresh token not found")
		}

		// For any other errors, return the error encountered
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	// Log the refresh token details for debugging (be mindful of sensitive information in production)
	log.Printf("Fetched refresh token: %+v", refreshToken)

	// Uncomment the token validation logic if necessary
	// if refreshToken.IsExpire || !refreshToken.IsActive {
	//     return nil, apperror.ErrInvalidToken
	// }

	return &refreshToken, nil
}

// RefreshToken implements RefreshTokenRepository.
func (r *refreshtokenRepository) RefreshToken(ctx context.Context, current_RT string, is_admin bool, env *appconfig.Env) (accesstoken string, refresh_token string, err error) {
	// Set a timeout for the refresh token process
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// Retrieve the current refresh token from the database
	re_token, err := r.GetRefreshTokenFromDB(ctx, current_RT, env)
	if err != nil {
		log.Printf("Error fetching refresh token from DB: %v", err)
		return "", "", fmt.Errorf("failed to retrieve refresh token: %w", err)
	}

	// Generate a new refresh token
	refresh_token, err = r.InsertRefreshTokenToDB(ctx, re_token.UserID, is_admin, env)
	if err != nil {
		log.Printf("Error inserting new refresh token into DB: %v", err)
		return "", "", fmt.Errorf("failed to insert new refresh token: %w", err)
	}

	// Update the old refresh token status to inactive
	_, err = r.UpdateRefreshTokenChanges(ctx, *re_token, env)
	if err != nil {
		log.Printf("Error updating old refresh token status: %v", err)
		return "", "", fmt.Errorf("failed to update old refresh token: %w", err)
	}

	// Create a new access token using the user ID
	id, err := primitive.ObjectIDFromHex(re_token.RefreshToken)
	if err != nil {
		log.Printf("Error converting refresh token to ObjectID: %v", err)
		return "", "", fmt.Errorf("failed to convert refresh token to ObjectID: %w", err)
	}

	// Generate the access token
	token, _, err := tokenutil.CreateAccessToken(id, env.ACCESS_SECRET, false, int(time.Second)*300)
	if err != nil {
		log.Printf("Error creating access token: %v", err)
		return "", "", fmt.Errorf("failed to create access token: %w", err)
	}

	log.Printf("Refresh token refreshed successfully for user ID: %s", re_token.UserID)
	return token, refresh_token, nil
}

// RevokeToken implements RefreshTokenRepository.
func (r *refreshtokenRepository) RevokeToken(ctx context.Context, current_RT string, env *appconfig.Env) error {
	// Set a timeout for the entire revoke operation
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	// Fetch the refresh token from the database
	RT, err := r.GetRefreshTokenFromDB(ctx, current_RT, env)
	if err != nil {
		log.Printf("Error retrieving refresh token for revocation: %v", err)
		return fmt.Errorf("failed to retrieve refresh token: %w", err)
	}

	// Log the details of the refresh token being revoked (avoid logging sensitive data in production)
	log.Printf("Revoking refresh token for token ID: %s", current_RT)

	// Revoke the token by updating the IsActive and IsExpire fields
	RT.RevokeAt = time.Now()
	RT.IsActive = false
	RT.IsExpire = true

	// Update the refresh token status in the database
	_, err = r.UpdateRefreshTokenChanges(ctx, *RT, env)
	if err != nil {
		log.Printf("Error updating refresh token status: %v", err)
		return fmt.Errorf("failed to update refresh token: %w", err)
	}

	log.Printf("Refresh token revoked successfully for token ID: %s", current_RT)
	return nil
}

// SaveRefreshTokenToDB implements RefreshTokenRepository.
func (r *refreshtokenRepository) InsertRefreshTokenToDB(ctx context.Context, userID string, is_admin bool, env *appconfig.Env) (string, error) {
	// Log function entry with userID
	log.Printf("Entering InsertRefreshTokenToDB for userID: %v", userID)

	// Convert userID string to ObjectID
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		// Log the error and return a wrapped error
		log.Printf("Error converting userID to ObjectID: %v, error: %v", userID, err)
		return "", fmt.Errorf("invalid user ID: %w", err)
	}

	// Generate the refresh token
	refreshToken, refresh_token_claims, err := tokenutil.CreateRefreshToken(oid, is_admin, env.REFRESH_SECRET, env.REFRESH_TOK_EXP)
	if err != nil {
		log.Printf("Error creating refresh token: %v", err)
		return "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	// Create the refresh token request object
	refreshTokenRequest := domain.RefreshTokenRequest{
		ID:           primitive.NewObjectID(),
		RefreshToken: refreshToken,
		UserID:       userID,
		ExpireAt:     refresh_token_claims.ExpiresAt.UTC().Local(), // Use the expiration value from config
		CreatedAt:    time.Now().UTC(),
		IsActive:     true,
		IsExpire:     false,
	}

	// Log the details of the token (you can remove this if not needed)
	log.Printf("Refresh token created for userID: %v, expires at: %v", userID, refreshTokenRequest.ExpireAt.Format("Mon Jan 2 15:04:05 2006"))

	// Insert the refresh token into the collection
	collection := r.db.Collection(r.collection_name)
	_, err = collection.InsertOne(ctx, refreshTokenRequest)
	if err != nil {
		log.Printf("Error inserting refresh token into DB for userID: %v, error: %v", userID, err)
		return "", fmt.Errorf("failed to insert refresh token into DB: %w", err)
	}

	// Log success
	log.Printf("Successfully inserted refresh token for userID: %v", userID)

	return refreshToken, nil
}

func NewrefreshTokenRepository(db *mongo.Database, collection_name string) RefreshTokenRepository {
	return &refreshtokenRepository{
		db:              db,
		collection_name: collection_name,
	}
}
