package repository

import (
	"context"
	"errors"
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
	RemoveItemFromCart(context context.Context, UserId primitive.ObjectID, product_id primitive.ObjectID, variant string) error
	AddProductItemIntoCart(context context.Context, UserID primitive.ObjectID, item domain.CartItem) (*domain.CartItem, error)
	UpdateProductItemByID(context context.Context, UserID primitive.ObjectID, updatedItem domain.CartItem) (*domain.CartItem, error)
}

type cartRepository struct {
	db              *mongo.Database
	collection_name string
}

// DONE
// AddProductItemIntoCart adds a product item into a user's cart by cart ID.
func (c *cartRepository) AddProductItemIntoCart(ctx context.Context, userID primitive.ObjectID, item domain.CartItem) (*domain.CartItem, error) {
	collection := c.db.Collection(c.collection_name)

	// Define the filter and update for the query
	filter := bson.M{"user_id": userID}
	update := bson.M{"$push": bson.M{"items": item}}

	// Perform the update operation
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to add item to cart: %w", err)
	}

	// Log the affected rows for debugging purposes
	if res.ModifiedCount == 0 {
		return nil, fmt.Errorf("no cart found with the specified ID: %s", userID.Hex())
	}

	return &item, nil
}

// 1. DONE
func (c *cartRepository) CreateCartByUserId(context context.Context, userID primitive.ObjectID) error {
	collection := c.db.Collection(c.collection_name)
	emptyCart := &domain.Cart{
		CartID:    primitive.NewObjectID(),
		UserID:    userID,
		Items:     make([]domain.CartItem, 0),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: nil,
	}
	_, err := collection.InsertOne(context, emptyCart)
	return err
}

// GetAllItemsInCartByUserID retrieves all items in a user's cart by their user_id.
func (r *cartRepository) GetAllItemsInCartByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.CartItem, error) {
	collection := r.db.Collection(r.collection_name)

	// Define the filter
	filter := bson.M{"user_id": userID}

	// Define a variable to store the result
	var cart domain.Cart

	// Query the database for the cart by user_id
	err := collection.FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Return an empty list and a custom error if no cart is found
			return nil, fmt.Errorf("no cart found for user with ID: %s", userID.Hex())
		}
		// Return a more generic error for unexpected issues
		return nil, fmt.Errorf("failed to fetch cart for user with ID %s: %w", userID.Hex(), err)
	}

	// Return the list of items from the cart
	return cart.Items, nil
}

// DONE
// GetCartByUserID retrieves the cart for a user by userID.
func (c *cartRepository) GetCartByUserID(ctx context.Context, userID primitive.ObjectID) (*domain.Cart, error) {
	collection := c.db.Collection(c.collection_name)

	// Define the filter for userID
	filter := bson.M{"user_id": userID}

	// Define the result variable
	var cart domain.Cart

	// Attempt to find the cart by userID
	err := collection.FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		// Handle different error scenarios
		if err == mongo.ErrNoDocuments {
			// Return a custom error for no cart found
			return nil, fmt.Errorf("no cart found for user with ID: %s", userID.Hex())
		}
		// Log the error with more details
		fmt.Printf("Error fetching cart for user %s: %v", userID.Hex(), err)
		// Wrap the error with additional context
		return nil, fmt.Errorf("failed to fetch cart for user ID %s: %w", userID.Hex(), err)
	}

	// Return the found cart
	return &cart, nil
}

// RemoveItemFromCart removes a product item from the cart.
func (r *cartRepository) RemoveItemFromCart(ctx context.Context, UserID primitive.ObjectID, productID primitive.ObjectID, variant string) error {
	collection := r.db.Collection(r.collection_name)

	// Define the filter to find the cart document
	cartFilter := bson.M{"user_id": UserID}

	// Define the update to pull the product item from the cart's items array
	update := bson.M{
		"$pull": bson.M{
			"items": bson.M{
				"product_id": productID,
				"variant":    variant,
			},
		},
		"$set": bson.M{
			"updatedAt": time.Now().UTC(),
		},
	}

	// Perform the update operation
	result, err := collection.UpdateOne(ctx, cartFilter, update)
	if err != nil {
		// Log the error with context
		fmt.Printf("Error removing item from cart (cartID: %s, productID: %s, variant: %s): %v", UserID.Hex(), productID.Hex(), variant, err)
		// Return a wrapped error with more context
		return err
	}

	// Check if any document was modified (item removed)
	if result.ModifiedCount == 0 {
		fmt.Printf("No items were removed from cart (cartID: %s, productID: %s, variant: %s)", UserID.Hex(), productID.Hex(), variant)
		// Optionally return an error or a specific message for no items modified
		return errors.New("Invalid product")
	}

	// Return nil if the operation was successful
	fmt.Printf("Item removed from cart (cartID: %s, productID: %s, variant: %s)", UserID.Hex(), productID.Hex(), variant)
	return nil
}

// DONE
func (c *cartRepository) UpdateProductItemByID(context context.Context, UserID primitive.ObjectID, updatedItem domain.CartItem) (*domain.CartItem, error) {
	collection := c.db.Collection(c.collection_name)
	// Define the filter to locate the specific cart and item
	filter := bson.M{
		"user_id": UserID,
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
			"items.$.discount":     updatedItem.Discount,
			"items.$.price":        updatedItem.Price,
			"items.$.buy_quantity": updatedItem.BuyQuantity, // Use the positional operator `$` to target the matched item
			"updatedAt":            time.Now().UTC(),
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
