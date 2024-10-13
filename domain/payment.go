package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type PayMethod int

const (
	Momo PayMethod = iota
	Bank
	Cash
)

type PayStatus int

const (
	Unpaid PayStatus = iota
	Paid
)

type Payment struct {
	PaymentID     primitive.ObjectID `json:"payment_id" bson:"payment_id"`
	PaymentMethod PayMethod          `json:"payment_method" bson:"payment_method"`
	PaymentStatus PayStatus          `json:"payment_status" bson:"payment_status"`
	InvoiceId     string             `json:"invoice_id" bson:"invoice_id"`
}
