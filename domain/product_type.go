package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductType struct {
	Product_type_id   primitive.ObjectID `json:"_id" bson:"_id"`
	Product_type_name string             `json:"type_name" bson:"type_name"`
	Description       string             `json:"description" bson:"description"`
}

type ProductTypeUsecase interface {
	CreateProductType(context context.Context, product_type ProductType) error
	GetAllProductType(context context.Context) ([]ProductType, error)
	GetProductTypeById(context context.Context, product_type_id int) (*ProductType, error)
	RemoveProductType(context context.Context, product_type_id int) error
	UpdateProductType(context context.Context, updated_product_type ProductType) (*ProductType, error)
}
