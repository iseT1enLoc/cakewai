package repository

import (
	"cakewai/cakewai.com/domain"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]*domain.User, error)
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userId string) error
}

type userRepository struct {
	db              *mongo.Database
	collection_name string
}

// DeleteUser implements UserRepository.
func (u *userRepository) DeleteUser(ctx context.Context, userId string) error {
	contx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	userObjectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatalf("Can not parsing the hex to object id")
		return err
	}

	_, err = u.db.Collection(u.collection_name).DeleteOne(contx, bson.M{"_id": userObjectID})
	return err
}

func NewUserRepository(db *mongo.Database, collection string) UserRepository {
	return &userRepository{
		db:              db,
		collection_name: collection,
	}
}

func (u *userRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// Set a default timeout for the operation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel() // Ensure the cancel function is called to release resources

	collection := u.db.Collection(u.collection_name)

	if user.GoogleId != nil {
		_, err := collection.InsertOne(ctx, bson.M{
			"_id":             user.Id,
			"google_id":       user.GoogleId,
			"name":            user.Name,
			"password":        user.Password,
			"email":           user.Email,
			"phone":           nil,
			"address":         nil,
			"profile_picture": user.ProfilePicture,
			"createdAt":       time.Now().UTC(),
			"updatedAt":       nil,
			"is_admin":        false,
		})

		if err != nil {
			log.Print(err)
			return nil, err
		}
		return user, nil
	}

	_, err := collection.InsertOne(ctx, bson.M{
		"_id":             user.Id,
		"google_id":       nil,
		"name":            user.Name,
		"password":        user.Password,
		"email":           user.Email,
		"phone":           nil,
		"address":         nil,
		"profile_picture": user.ProfilePicture,
		"createdAt":       time.Now().UTC(),
		"updatedAt":       nil,
		"is_admin":        false,
	})

	if err != nil {
		log.Print(err)
		return nil, err
	}
	return user, nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	// Set a default timeout for the operation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel() // Ensure resources are released

	collection := u.db.Collection(u.collection_name)
	var user domain.User

	// Perform the query
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Handle case where no user is found
			log.Printf("No user found with email: %s", email)
			log.Println(err)
			return nil, err // Return nil user and nil error
		}
		// Handle other errors (e.g., timeout)
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Timeout occurred while finding user with email: %s", email)
		} else {
			log.Printf("Failed to find user with email: %s, error: %v", email, err)
		}
		return nil, err
	}

	return &user, nil
}

// GetUserById implements UserRepository.
func (u *userRepository) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	// Set a timeout for the operation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel() // Ensure resources are released

	collection := u.db.Collection(u.collection_name)

	// Convert the string ID to a MongoDB ObjectID
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %s, error: %v", id, err)
		return nil, err
	}

	// Define a variable to hold the result
	var fuser domain.User

	// Query the database
	err = collection.FindOne(ctx, bson.M{"_id": ObjectID}).Decode(&fuser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("No user found with ID: %s", id)
			return nil, nil // Return nil for both user and error if no user is found
		}
		log.Printf("Failed to fetch user with ID: %s, error: %v", id, err)
		return nil, err
	}

	// Return the found user
	log.Printf("User retrieved successfully: %s", fuser.Name)
	return &fuser, nil
}

// GetUsers implements UserRepository.
func (u *userRepository) GetUsers(ctx context.Context) ([]*domain.User, error) {
	// Set a timeout for the operation
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel() // Ensure resources are released

	collection := u.db.Collection(u.collection_name)

	// Exclude the "password" field from the results for security reasons
	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})

	// Execute the query
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Printf("Failed to retrieve users: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx) // Ensure the cursor is closed to avoid resource leaks

	// Parse the results into a slice of users
	var users []*domain.User
	if err = cursor.All(ctx, &users); err != nil {
		log.Printf("Error decoding users: %v", err)
		return nil, err
	}

	// Handle the case where no users are found
	if len(users) == 0 {
		log.Println("No users found in the database")
		return []*domain.User{}, nil // Return an empty slice and no error
	}

	log.Printf("Retrieved %d users from the database", len(users))
	return users, nil
}

// UpdateUser implements UserRepository.
func (u *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	// Set a timeout for the operation
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel() // Ensure resources are released

	collection := u.db.Collection(u.collection_name)

	// Construct the update fields
	updateFields := bson.M{
		"$set": bson.M{
			"google_id":       user.GoogleId,
			"profile_picture": user.ProfilePicture,
			"name":            user.Name,
			"password":        user.Password,
			"email":           user.Email,
			"phone":           user.Phone,
			"address":         user.Address,
			"created_at":      user.CreatedAt,
			"updated_at":      time.Now(), // Update `updated_at` to the current time
			"role_id":         user.RoleID,
		},
	}

	// Perform the update operation
	result, err := collection.UpdateOne(ctx, bson.M{"_id": user.Id}, updateFields)
	if err != nil {
		log.Printf("Failed to update user %v: %v", user.Id, err)
		return err
	}

	// Check if any documents were updated
	if result.MatchedCount == 0 {
		log.Printf("No user found with the provided ID: %v", user.Id)
		return fmt.Errorf("no user found with the provided ID")
	}

	// Log a successful update
	log.Printf("User with ID %v updated successfully", user.Id)

	return nil
}
