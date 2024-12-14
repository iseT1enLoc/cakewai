package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	PaymentMethod string `json:"payment_method" bson:"payment_method" default:"cash"`
	IsPaid        int    `json:"is_paid" bson:"is_paid" default:"false"`
}
type Order struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	CustomerID      primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	CustomerName    string             `json:"customer_name" bson:"customer_name"`
	ShippingAddress string             `json:"shipping_address" bson:"shipping_address"`
	PhoneNumber     string             `json:"phone_number" bson:"phone_number"`
	Notes           string             `json:"notes" bson:"notes"`
	OrderItems      []CartItem         `json:"order_items" bson:"order_items"`
	TotalPrice      float64            `json:"total_price" bson:"total_price"`

	PaymentInfo    Payment    `json:"payment_info" bson:"payment_info"`
	OrderStatus    string     `json:"order_status" bson:"order_status"  default:"pending"` // Can be "pending", "shipped", "delivered", "canceled"
	ShippingStatus string     `json:"shipping_status" bson:"shipping_status" default:"pending"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt      *time.Time `json:"update_at" bson:"updated_at"`
}
type PaymentReq struct {
	Order_id    string  `json:"order_id" bson:"order_id"`
	PaymentInfo Payment `json:"payment_info" bson:"payment_info"`
}
type OrderReq struct {
	CustomerName    string     `json:"customer_name" bson:"customer_name"`
	ShippingAddress string     `json:"shipping_address" bson:"shipping_address"`
	PhoneNumber     string     `json:"phone_number" bson:"phone_number"`
	Notes           string     `json:"notes" bson:"notes"`
	OrderItems      []CartItem `json:"order_items" bson:"order_items"`
}
type OrderUsecase interface {
	CreateOrder(context context.Context, order Order) (*Order, error)
	GetAllOrders(context context.Context) ([]*Order, error)
	UpdateOrder(context context.Context, updatedOrder Order) (*Order, error)
	UpdateOrderPaymentStatus(context context.Context, order_id primitive.ObjectID, is_paid int) (int, error)
	GetOrderByID(context context.Context, ID primitive.ObjectID) (*Order, error)
}
