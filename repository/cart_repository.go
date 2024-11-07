package repository

import (
	"context"
	"fmt"
	"time"

	"cakewai/cakewai.com/domain"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CartRepository interface {
	CreateCartByUserId(context context.Context, UserId primitive.ObjectID) error
	GetAllItemsInCartByUserID(context context.Context, user_id primitive.ObjectID) ([]domain.CartItem, error)
	GetCartByUserID(context context.Context, UserId primitive.ObjectID) (*domain.Cart, error)
	RemoveItemFromCart(context context.Context, cartID primitive.ObjectID, product_id primitive.ObjectID, variant string) error
	AddProductItemIntoCart(context context.Context, cardid primitive.ObjectID, item domain.CartItem) (*domain.CartItem, error)
	UpdateProductItemByID(context context.Context, cardID primitive.ObjectID, updatedItem domain.CartItem) (*domain.CartItem, error)
}
type cartRepository struct {
	db              *mongo.Database
	collection_name string
}

// DONE
func (c *cartRepository) AddProductItemIntoCart(context context.Context, cartid primitive.ObjectID, item domain.CartItem) (*domain.CartItem, error) {
	collection := c.db.Collection(c.collection_name)

	update := bson.D{{"$push", bson.M{"items": item}}} //push||set||... read later
	filter := bson.M{"_id": cartid}
	res, err := collection.UpdateOne(context, filter, update)
	fmt.Printf("Row EFFECT %v", res.ModifiedCount)
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
func (c *cartRepository) GetAllItemsInCartByUserID(context context.Context, user_id primitive.ObjectID) ([]domain.CartItem, error) {
	collection := c.db.Collection(c.collection_name)
	var card domain.Cart
	err := collection.FindOne(context, bson.M{"user_id": user_id}).Decode(&card)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return card.Items, nil
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
	// Define the filter to find the document
	cart_item_filter := bson.M{"_id": cartID}

	// Define the update to pull an item from the items array
	update := bson.M{
		"$pull": bson.M{
			"items": bson.M{
				"product_id": product_id,
				"variant":    variant,
			},
		},
	}
	_, err := collection.UpdateOne(context, cart_item_filter, update)

	log.Error(err)
	return err
}

// DONE
func (c *cartRepository) UpdateProductItemByID(context context.Context, cartID primitive.ObjectID, updatedItem domain.CartItem) (*domain.CartItem, error) {
	collection := c.db.Collection(c.collection_name)
	// Define the filter to locate the specific cart and item
	filter := bson.M{
		"_id": cartID,
		"items": bson.M{
			"$elemMatch": bson.M{
				"product_id": updatedItem.ProductId,
				"variant":    updatedItem.Variant,
			},
		},
	}

	// Define the update to set a new quantity
	update := bson.M{
		"$set": bson.M{
			"items.$.buy_quantity": updatedItem.BuyQuantity, // Use the positional operator `$` to target the matched item
		},
	}
	res, err := collection.UpdateOne(context, filter, update)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if res.ModifiedCount == 0 {
		fmt.Println("No matching item found to update")
	} else {
		fmt.Println("Item updated successfully!")
	}
	return &updatedItem, nil
}

func NewCartRepository(db *mongo.Database, collection_name string) CartRepository {
	return &cartRepository{
		db:              db,
		collection_name: collection_name,
	}
}
