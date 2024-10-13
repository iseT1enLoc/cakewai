package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	InvoiceID  primitive.ObjectID `json:"invoice_id" bson:"invoice_id"`
	PaymentID  string             `json:"payment_id" bson:"payment_id"`
	CustomerID string             `json:"customer_id" bson:"customer_id"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	TotalPrice float64            `json:"total_price" bson:"total_price"`
}
