package repository

import (
	"cakewai/cakewai.com/domain"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"

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
	UpdateUserPassword(ctx context.Context, userId primitive.ObjectID, currentPassword string, newPassword string) error
	HandleForgotPassword(ctx context.Context, email string) (*string, *string, error)
}

type userRepository struct {
	db              *mongo.Database
	collection_name string
}

// HandleForgotPassword implements UserRepository.
func (u *userRepository) HandleForgotPassword(ctx context.Context, email string) (*string, *string, error) {
	// Set a timeout for the operation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel() // Ensure resources are released

	var fuser domain.User
	collection := u.db.Collection(u.collection_name)

	// Query the database to find the user by email
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&fuser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("No user found with email: %s", email)
			return nil, nil, fmt.Errorf("no user found with email: %s", email)
		}
		log.Printf("Failed to fetch user with Email: %s, error: %v", email, err)
		return nil, nil, err
	}

	// Generate a random 8-digit number (for the new password)
	rand.Seed(uint64(time.Now().Hour())) // Proper seeding for randomness
	randomNumber := rand.Intn(90000000) + 10000000

	// Encrypt the random number to use as the new password
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("%d", randomNumber)), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return nil, nil, err
	}

	// Update the user's password in the database
	updateFields := bson.M{
		"$set": bson.M{
			"password": encryptedPassword,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"email": email}, updateFields)
	if err != nil {
		log.Printf("Failed to update password for email: %s, error: %v", email, err)
		return nil, nil, err
	}

	// Return the generated random password as a string
	result := fmt.Sprintf("%d", randomNumber)
	return &fuser.Name, &result, nil
}

// UpdateUserPassword implements UserRepository.
func (u *userRepository) UpdateUserPassword(ctx context.Context, userId primitive.ObjectID, currentPassword string, newPassword string) error {
	// Set a timeout for the operation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel() // Ensure resources are released

	collection := u.db.Collection(u.collection_name)

	// Convert the string ID to a MongoDB ObjectID
	ObjectID, err := primitive.ObjectIDFromHex(userId.Hex())
	if err != nil {
		log.Printf("Invalid ID format: %s, error: %v", userId, err)
		return err
	}

	// Define a variable to hold the result
	var fuser domain.User

	// Query the database
	err = collection.FindOne(ctx, bson.M{"_id": ObjectID}).Decode(&fuser)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("No user found with ID: %s", userId)
			return err
		}
		log.Printf("Failed to fetch user with ID: %s, error: %v", userId, err)
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(fuser.Password), []byte(currentPassword)); err != nil {

		err = errors.New("Invalid current password")
		return err
	}
	encryptednewPassword, err := bcrypt.GenerateFromPassword(
		[]byte(newPassword),
		bcrypt.DefaultCost,
	)
	fuser.Password = string(encryptednewPassword)

	updateFields := bson.M{
		"$set": bson.M{
			"password": encryptednewPassword,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": userId}, updateFields)
	return err
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
			"address":         domain.Address{},
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
		"address":         domain.Address{},
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
	fmt.Println("Print at line 152 and below is user")
	fmt.Println(fuser)
	fmt.Printf("\nUser homecode %s\n", fuser.Address.Homecode)
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
	// ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel() // Ensure resources are released

	collection := u.db.Collection(u.collection_name)
	fmt.Println(user.Password)
	fmt.Println(user)
	fmt.Println("enter line 210")
	var updateFields bson.M
	if user.Password != "" {
		encryptedPassword, _ := bcrypt.GenerateFromPassword(
			[]byte(user.Password),
			bcrypt.DefaultCost,
		)
		// Construct the update fields
		updateFields = bson.M{
			"$set": bson.M{
				"google_id":            user.GoogleId,
				"profile_picture":      user.ProfilePicture,
				"name":                 user.Name,
				"password":             encryptedPassword,
				"email":                user.Email,
				"phone":                user.Phone,
				"created_at":           user.CreatedAt,
				"updated_at":           time.Now(), // Update `updated_at` to the current time
				"address.home_code":    user.Address.Homecode,
				"address.street":       user.Address.Street,
				"address.district":     user.Address.District,
				"address.province":     user.Address.Province,
				"address.full_address": user.Address.FullAddress,
			},
		}
	} else {
		// Construct the update fields
		updateFields = bson.M{
			"$set": bson.M{
				"google_id":            user.GoogleId,
				"profile_picture":      user.ProfilePicture,
				"name":                 user.Name,
				"email":                user.Email,
				"phone":                user.Phone,
				"created_at":           user.CreatedAt,
				"updated_at":           time.Now(), // Update `updated_at` to the current time
				"address.home_code":    user.Address.Homecode,
				"address.street":       user.Address.Street,
				"address.district":     user.Address.District,
				"address.province":     user.Address.Province,
				"address.full_address": user.Address.FullAddress,
			},
		}
	}

	fmt.Println(user.Id)
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
