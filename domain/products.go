package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ProductId     primitive.ObjectID `json:"_id" bson:"_id"`
	ProductName   string             `json:"product_name" bson:"product_name"`
	Slug          string             `json:"slug"`
	Description   string             `json:"description" bson:"description"`
	ImageLink     string             `json:"image_link" bson:"image_link"`
	ProductTypeID string             `json:"product_type_id" bson:"product_type_id"`
	Variant       []ProductVariant   `json:"product_variant" bson:"product_variant"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
}
type ProductVariant struct {
	VarientFeatures string  `json:"variant_features" bson:"variant_features"`
	Price           float64 `json:"price" bson:"price"`
	NumInStore      int     `json:"num_in_store" bson:"num_in_store" default:"100"`
	Discount        float64 `json:"discount" bson:"discount"`
}
type ProductRequest struct {
	Description   string           `json:"description" bson:"description"`
	ProductName   string           `json:"product_name" bson:"product_name"`
	ImageLink     string           `json:"image_link" bson:"image_link"`
	ProductTypeID string           `json:"product_type_id" bson:"product_type_id"`
	Variant       []ProductVariant `json:"product_variant" bson:"product_variant"`
}

type ProductUsecase interface {
	// Create a new product
	CreateProduct(ctx context.Context, product *ProductRequest) (*Product, error)

	// Get product by ID
	GetProductById(ctx context.Context, id primitive.ObjectID) (*Product, error)

	GetProductBySlug(ctx context.Context, slug string) (*Product, error)
	// Get all products
	GetAllProducts(ctx context.Context) ([]*Product, error)

	// Update product by ID
	UpdateProductById(ctx context.Context, id primitive.ObjectID, updatedProduct *ProductRequest) (int64, error)

	// Delete product by ID
	DeleteProductById(ctx context.Context, id primitive.ObjectID) error

	GetProductByProductTypeID(ctx context.Context, id string) ([]*Product, error)
	// Add a variant to a product
	AddProductVariant(ctx context.Context, productId primitive.ObjectID, variant ProductVariant) (int64, error)

	// Update a product variant
	UpdateProductVariant(ctx context.Context, productId primitive.ObjectID, updatedVariant ProductVariant) (int64, error)

	// Delete a product variant
	DeleteProductVariant(ctx context.Context, productId primitive.ObjectID, variant_name string) (int64, error)

	//Search cake
	SearchProducts(ctx context.Context, query string, type_id string, variant string) ([]*Product, error)

	FetchSortedProducts(ctx context.Context, sortField string, sortOrder string) ([]*Product, error)
}
