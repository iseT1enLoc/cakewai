package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
	ProductId   primitive.ObjectID `json:"product_id,omitempty" bson:"product_id,omitempty"`
	ProductType string             `json:"type_id" bson:"type_id"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Variant     string             `json:"variant" bson:"variant"`
	Discount    float64            `json:"discount" bson:"discount"`
	Price       float64            `json:"price,omitempty" bson:"price,omitempty"`
	ImageLink   string             `json:"image_link" bson:"image_link"`
	BuyQuantity int                `json:"buy_quantity" bson:"buy_quantity"`
}
type Cart struct {
	CartID    primitive.ObjectID `json:"_id" bson:"_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Items     []CartItem         `json:"items" bson:"items"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt *time.Time         `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type CartUsecase interface {
	CreateCartByUserId(context context.Context, UserId primitive.ObjectID) error

	GetAllItemsInCartByUserID(context context.Context, user_id primitive.ObjectID) ([]CartItem, error)

	GetCartByUserID(context context.Context, userID primitive.ObjectID) (*Cart, error)

	RemoveItemFromCart(context context.Context, userID primitive.ObjectID, product_id primitive.ObjectID, variant string) error

	AddCartItemIntoCart(context context.Context, userID primitive.ObjectID, item CartItem) (*primitive.ObjectID, error)

	UpdateAnCartItemByUserID(context context.Context, userID primitive.ObjectID, updatedItem CartItem) (*CartItem, error)
}
