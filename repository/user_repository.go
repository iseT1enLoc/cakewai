package repository

import (
	"cakewai/cakewai.com/domain"
	"context"
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
	collection := u.db.Collection(u.collection_name)

	if user.GoogleId != "" {
		_, err := collection.InsertOne(ctx, bson.M{
			"_id":             user.Id,
			"email":           user.Email,
			"google_id":       user.GoogleId,
			"name":            user.Name,
			"profile_picture": user.ProfilePicture,
			"password":        user.Password,
			"createdAt":       time.Now(),
		})

		if err != nil {
			log.Print(err)
			return nil, err
		}
		return user, nil
	}

	_, err := collection.InsertOne(ctx, bson.M{
		"_id":             user.Id,
		"email":           user.Email,
		"google_id":       user.GoogleId,
		"name":            user.Name,
		"password":        user.Password,
		"profile_picture": user.ProfilePicture,
		"createdAt":       time.Now(),
	})

	if err != nil {
		log.Print(err)
		return nil, err
	}
	return user, nil
}

// GetUserByEmail implements UserRepository.
func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	collection := u.db.Collection(u.collection_name)
	var user domain.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserById implements UserRepository.
func (u *userRepository) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	collection := u.db.Collection(u.collection_name)

	var fuser domain.User
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Print(err)
	}
	err = collection.FindOne(ctx, bson.M{"_id": ObjectID}).Decode(&fuser)
	print(fuser.Name)
	return &fuser, nil
}

// GetUsers implements UserRepository.
func (u *userRepository) GetUsers(ctx context.Context) ([]*domain.User, error) {
	collection := u.db.Collection(u.collection_name)

	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collection.Find(ctx, bson.D{}, opts)

	if err != nil {
		return nil, err
	}

	var users []*domain.User

	err = cursor.All(ctx, &users)
	if users == nil {
		return []*domain.User{}, err
	}

	return users, err
}

// UpdateUser implements UserRepository.
func (u *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	collection := u.db.Collection(u.collection_name)

	_, err := collection.UpdateOne(ctx, bson.M{"_id": user.Id},
		bson.M{"$set": bson.M{"google_id": user.GoogleId,
			"profile_picture": user.ProfilePicture,
			"name":            user.Name,
			"password":        user.Password,
			"email":           user.Email,
			"phone":           user.Phone,
			"address":         user.Address,
			"created_at":      user.CreatedAt,
			"updated_at":      user.UpdatedAt,
			"role_id":         user.RoleID,
		}})

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
