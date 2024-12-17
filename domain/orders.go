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
	CustomerName    string             `json:"name" bson:"name"`
	ShippingAddress string             `json:"shipping_address" bson:"shipping_address"`
	PhoneNumber     string             `json:"phone" bson:"phone"`
	Notes           string             `json:"notes" bson:"notes"`
	Email           string             `json:"email" bson:"email"`
	OrderItems      []CartItem         `json:"order_items" bson:"order_items"`
	TotalPrice      float64            `json:"total_price" bson:"total_price"`

	ServiceType    int       `json:"service_type" bson:"service_type"` //1-wwebsite, 2-AI
	PaymentInfo    Payment   `json:"payment_info" bson:"payment_info"`
	OrderStatus    string    `json:"order_status" bson:"order_status"  default:"pending"` // Can be "pending", "shipped", "delivered", "canceled"
	ShippingStatus string    `json:"shipping_status" bson:"shipping_status" default:"pending"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"update_at" bson:"updated_at"`
}
type PaymentReq struct {
	PaymentInfo Payment `json:"payment_info" bson:"payment_info"`
}
type OrderReq struct {
	CustomerName    string     `json:"name" bson:"name"`
	ShippingAddress string     `json:"address" bson:"address"`
	PhoneNumber     string     `json:"phone" bson:"phone"`
	Notes           string     `json:"notes" bson:"notes"`
	Email           string     `json:"email" bson:"email"`
	OrderItems      []CartItem `json:"order_items" bson:"order_items"`
	ServiceType     int        `json:"service_type" bson:"service_type"`
}
type OrderUsecase interface {
	CreateOrder(context context.Context, order Order) (*Order, error)
	GetAllOrders(context context.Context) ([]*Order, error)
	UpdateOrder(context context.Context, updatedOrder Order) (*Order, error)
	UpdateOrderPaymentStatus(context context.Context, order_id primitive.ObjectID, is_paid int) (int, error)
	GetOrderByID(context context.Context, ID primitive.ObjectID) (*Order, error)
	GetOrdersByCustomerID(context context.Context, CustomerID primitive.ObjectID) ([]*Order, error)
	DeleteOrder(context context.Context, order_id primitive.ObjectID) error
}
