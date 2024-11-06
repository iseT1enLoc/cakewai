package repository

import (
	"context"
	"time"

	"cakewai/cakewai.com/domain"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartRepository interface {
	CreateCartByUserId(context context.Context, UserId primitive.ObjectID) error
	GetAllItemsInCartByUserID(context context.Context, user_id primitive.ObjectID) ([]*domain.CartItem, error)
	GetCartByUserID(context context.Context, UserId primitive.ObjectID) (*domain.Cart, error)
	RemoveItemFromCart(context context.Context, cartID primitive.ObjectID, product_id primitive.ObjectID, variant string) error
	AddProductItemIntoCart(context context.Context, item domain.CartItem) (*domain.CartItem, error)
	UpdateProductItemByID(context context.Context, updatedItem domain.CartItem) (*domain.CartItem, error)
}
type cartRepository struct {
	db              *mongo.Database
	collection_name string
}

// DONE
func (c *cartRepository) AddProductItemIntoCart(context context.Context, item domain.CartItem) (*domain.CartItem, error) {
	collection := c.db.Collection(c.collection_name)
	_, err := collection.InsertOne(context, item)
	if err != nil {
		return nil, err
	}
	return &item, err
}

// DONE
func (c *cartRepository) CreateCartByUserId(context context.Context, userID primitive.ObjectID) error {
	collection := c.db.Collection(c.collection_name)
	emptyCart := &domain.Cart{
		CartID:    primitive.NewObjectID(),
		UserID:    userID,
		Items:     make([]domain.CartItem, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := collection.InsertOne(context, emptyCart)
	return err
}

// DONE
func (c *cartRepository) GetAllItemsInCartByUserID(context context.Context, user_id primitive.ObjectID) ([]*domain.CartItem, error) {
	collection := c.db.Collection(c.collection_name)
	items_cursor, err := collection.Find(context, bson.D{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var items_list []*domain.CartItem
	err = items_cursor.All(context, &items_list)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return items_list, nil
}

// DONE
func (c *cartRepository) GetCartByUserID(context context.Context, userID primitive.ObjectID) (*domain.Cart, error) {
	collection := c.db.Collection(c.collection_name)
	cart_filter := bson.M{"user_id": userID}
	var cart *domain.Cart
	err := collection.FindOne(context, cart_filter).Decode(&cart)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return cart, nil
}

// DONE
func (c *cartRepository) RemoveItemFromCart(context context.Context, cartID primitive.ObjectID, product_id primitive.ObjectID, variant string) error {
	collection := c.db.Collection(c.collection_name)
	cart_item_filter := bson.M{"product_id": product_id, "variant": variant}
	_, err := collection.DeleteOne(context, cart_item_filter)
	log.Error(err)
	return err
}

// DONE
func (c *cartRepository) UpdateProductItemByID(context context.Context, updatedItem domain.CartItem) (*domain.CartItem, error) {
	collection := c.db.Collection(c.collection_name)
	cart_item_filter := bson.M{"product_id": updatedItem.ProductId, "variant": updatedItem.Variant}
	_, err := collection.UpdateOne(context, cart_item_filter, updatedItem)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &updatedItem, nil
}

func NewCartRepository(db *mongo.Database, collection_name string) CartRepository {
	return &cartRepository{
		db:              db,
		collection_name: collection_name,
	}
}
