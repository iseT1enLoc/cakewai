package usecase

import (
	"cakewai/cakewai.com/domain"
	"cakewai/cakewai.com/repository"
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type cartUsecase struct {
	cartRepository repository.CartRepository
	timeout        time.Duration
}

// AddCartItemIntoCart implements domain.CartUsecase.
func (c *cartUsecase) AddCartItemIntoCart(context context.Context, userID primitive.ObjectID, item domain.CartItem) (*primitive.ObjectID, error) {
	// Check if the user's cart exists
	cart, err := c.cartRepository.GetCartByUserID(context, userID)
	if err != nil {
		if strings.Contains(err.Error(), "mongo: no documents in result") {
			return nil, fmt.Errorf("cart not found for user: %s", userID.Hex())
		}
		return nil, fmt.Errorf("failed to fetch cart for user %s: %w", userID.Hex(), err)
	}

	// Check if the item already exists in the cart
	for _, cartItem := range cart.Items {
		if cartItem.ProductId == item.ProductId && cartItem.Variant == item.Variant {
			return nil, fmt.Errorf("item with product_id %s and variant %s already exists in the cart", item.ProductId.Hex(), item.Variant)
		}
	}

	// Add the item to the cart
	addedItem, err := c.cartRepository.AddProductItemIntoCart(context, cart.CartID, item)
	if err != nil {
		return nil, fmt.Errorf("failed to add item to cart for user %s: %w", userID.Hex(), err)
	}

	// Return the ID of the added item
	return &addedItem.ProductId, nil
}

// CreateCartByUserId implements domain.CartUsecase.
func (c *cartUsecase) CreateCartByUserId(context context.Context, UserId primitive.ObjectID) error {

	// Create an empty cart using the repository method
	err := c.cartRepository.CreateCartByUserId(context, UserId)
	if err != nil {
		return fmt.Errorf("failed to create cart for user %s: %w", UserId.Hex(), err)
	}

	return nil
}

// GetAllItemsInCartByUserID implements domain.CartUsecase.
func (c *cartUsecase) GetAllItemsInCartByUserID(context context.Context, user_id primitive.ObjectID) ([]domain.CartItem, error) {
	// Call the repository to get the items in the cart
	items, err := c.cartRepository.GetAllItemsInCartByUserID(context, user_id)
	if err != nil {
		// Handle the case where no cart or items exist for the user
		if strings.Contains(err.Error(), "mongo: no documents in result") {
			return nil, fmt.Errorf("no cart found for user: %s", user_id.Hex())
		}

		// Generic error
		return nil, fmt.Errorf("failed to fetch cart items: %w", err)
	}

	// Return the list of items
	return items, nil
}

// GetCartByUserID implements domain.CartUsecase.
func (c *cartUsecase) GetCartByUserID(context context.Context, userID primitive.ObjectID) (*domain.Cart, error) {
	// Call the repository to fetch the cart
	cart, err := c.cartRepository.GetCartByUserID(context, userID)
	if err != nil {
		// Handle case where cart is not found
		if strings.Contains(err.Error(), "mongo: no documents in result") {
			return nil, fmt.Errorf("cart not found for user: %s", userID.Hex())
		}

		// Generic error
		return nil, fmt.Errorf("failed to fetch cart: %w", err)
	}

	// Return the retrieved cart
	return cart, nil
}

// RemoveItemFromCart implements domain.CartUsecase.
func (c *cartUsecase) RemoveItemFromCart(context context.Context, userID primitive.ObjectID, product_id primitive.ObjectID, variant string) error {
	// Call the repository to remove the item
	err := c.cartRepository.RemoveItemFromCart(context, userID, product_id, variant)
	if err != nil {
		// Check if the item was not found in the cart
		if strings.Contains(err.Error(), "no matching item found") {
			return fmt.Errorf("cart item not found for user: %s", userID.Hex())
		}

		// Generic error
		return fmt.Errorf("failed to remove cart item: %w", err)
	}

	return nil
}

// UpdateAnCartItemByUserID implements domain.CartUsecase.
func (c *cartUsecase) UpdateAnCartItemByUserID(context context.Context, userID primitive.ObjectID, updatedItem domain.CartItem) (*domain.CartItem, error) {

	// Call the repository to update the item
	updatedCartItem, err := c.cartRepository.UpdateProductItemByID(context, userID, updatedItem)
	if err != nil {
		// Handle "no matching item" error specifically
		if strings.Contains(err.Error(), "no matching cart item found") {
			return nil, fmt.Errorf("cart item not found for user: %s", userID.Hex())
		}

		// Generic error
		return nil, fmt.Errorf("failed to update cart item: %w", err)
	}

	return updatedCartItem, nil
}

func NewCartUsecase(cart_repo repository.CartRepository, timeout time.Duration) domain.CartUsecase {
	return &cartUsecase{
		cartRepository: cart_repo,
		timeout:        timeout,
	}
}
