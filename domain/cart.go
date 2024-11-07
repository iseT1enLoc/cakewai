package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
	ProductId   primitive.ObjectID `json:"product_id,omitempty" bson:"product_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Variant     string             `json:"variant" bson:"variant"`
	ImageLink   string             `json:"image_link" bson:"image_link"`
	Price       float64            `json:"price,omitempty" bson:"price,omitempty"`
	BuyQuantity int                `json:"buy_quantity" bson:"buy_quantity"`
}
type Cart struct {
	CartID    primitive.ObjectID `json:"cart_id" bson:"_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Items     []CartItem         `json:"items" bson:"items"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type CartUsecase interface {
	CreateCartByUserId(context context.Context, UserId primitive.ObjectID) error
	GetAllItemsInCartByUserID(context context.Context, user_id primitive.ObjectID) ([]CartItem, error)
	GetCartByUserID(context context.Context, userID primitive.ObjectID) (*Cart, error)
	RemoveItemFromCart(context context.Context, cartID primitive.ObjectID, product_id primitive.ObjectID, variant string) error
	AddCartItemIntoCart(context context.Context, cardid primitive.ObjectID, item CartItem) (*primitive.ObjectID, error)
	UpdateCartItemByID(context context.Context, carID primitive.ObjectID, updatedItem CartItem) (*CartItem, error)
}
