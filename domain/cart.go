package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
	ProductId   primitive.ObjectID `json:"productId,omitempty" bson:"productId,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Price       float64            `json:"price,omitempty" bson:"price,omitempty"`
	BuyQuantity int                `json:"buyQuantity" bson:"buyQuantity"`
	CreatedAt   time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
