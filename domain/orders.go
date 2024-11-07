package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	PaymentMethod string `json:"payment_method" bson:"payment_method"`
	IsPaid        bool   `json:"is_paid" bson:"is_paid"`
}
type Order struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	CustomerID      primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	OrderItems      []CartItem         `json:"order_items" bson:"order_items"`
	TotalPrice      float64            `json:"total_price" bson:"total_price"`
	OrderStatus     string             `json:"order_status" bson:"order_status"` // Can be "pending", "shipped", "delivered", "canceled"
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"update_at" bson:"updated_at"`
	ShippingAddress Address            `json:"shipping_address" bson:"shipping_address"`
	PaymentInfo     Payment            `json:"payment_info" bson:"payment_info"`
}
type PaymentReq struct {
	Order_id string `json:"order_id"`
	Is_paid  int    `json:"is_paid"`
}
type OrderUsecase interface {
	CreateOrder(context context.Context, order Order) (*Order, error)
	GetAllOrders(context context.Context) ([]*Order, error)
	UpdateOrder(context context.Context, updatedOrder Order) (*Order, error)
	UpdateOrderPaymentStatus(context context.Context, order_id primitive.ObjectID, is_paid int) (int, error)
	GetOrderByID(context context.Context, ID primitive.ObjectID) (*Order, error)
}
