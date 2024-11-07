package usecase

import (
	"cakewai/cakewai.com/domain"
	"cakewai/cakewai.com/repository"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type cartUsecase struct {
	cartRepository repository.CartRepository
	timeout        time.Duration
}

// DONE
func (c *cartUsecase) AddCartItemIntoCart(context context.Context, cardid primitive.ObjectID, item domain.CartItem) (*primitive.ObjectID, error) {
	_, err := c.cartRepository.AddProductItemIntoCart(context, cardid, item)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &item.ProductId, nil
}

// DONE
func (c *cartUsecase) CreateCartByUserId(context context.Context, UserId primitive.ObjectID) error {
	err := c.cartRepository.CreateCartByUserId(context, UserId)
	return err
}

// DONE
func (c *cartUsecase) GetAllItemsInCartByUserID(context context.Context, user_id primitive.ObjectID) ([]domain.CartItem, error) {
	items, err := c.cartRepository.GetAllItemsInCartByUserID(context, user_id)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return items, err
}

// DONE
func (c *cartUsecase) GetCartByUserID(context context.Context, userID primitive.ObjectID) (*domain.Cart, error) {
	cart, err := c.cartRepository.GetCartByUserID(context, userID)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return cart, nil
}

// DONE
func (c *cartUsecase) RemoveItemFromCart(context context.Context, cartID primitive.ObjectID, product_id primitive.ObjectID, variant string) error {
	err := c.cartRepository.RemoveItemFromCart(context, cartID, product_id, variant)
	log.Print(err)
	return err
}

// DONE
func (c *cartUsecase) UpdateCartItemByID(context context.Context, cardID primitive.ObjectID, updatedItem domain.CartItem) (*domain.CartItem, error) {
	cart_item, err := c.cartRepository.UpdateProductItemByID(context, cardID, updatedItem)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return cart_item, nil
}

func NewCartUsecase(cart_repo repository.CartRepository, timeout time.Duration) domain.CartUsecase {
	return &cartUsecase{
		cartRepository: cart_repo,
		timeout:        timeout,
	}
}
