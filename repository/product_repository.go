package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cakewai/cakewai.com/domain"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	// Create a new product: DONE-------|DONE handler
	CreateProduct(ctx context.Context, product *domain.ProductRequest) (*domain.Product, error)

	// Get product by ID:DONE------------|DONE handler
	GetProductById(ctx context.Context, id primitive.ObjectID) (*domain.Product, error)

	// Get all products: DONE------------|DONE handler
	GetAllProducts(ctx context.Context) ([]*domain.Product, error)

	// Update product by ID:DONE----------|DONE handler
	UpdateProductById(ctx context.Context, id primitive.ObjectID, updatedProduct *domain.ProductRequest) (int64, error)

	// Delete product by ID: DONE-------|DONE handler
	DeleteProductById(ctx context.Context, id primitive.ObjectID) (rowAffect int64, err error)

	// Add a variant to a product:DONE----------|DONE handler
	AddProductVariant(ctx context.Context, productId primitive.ObjectID, variant domain.ProductVariant) (int64, error)

	// Update a product variant
	UpdateProductVariant(ctx context.Context, productId primitive.ObjectID, updatedVariant domain.ProductVariant) (int64, error)

	// Delete a product variant
	DeleteProductVariant(ctx context.Context, productId primitive.ObjectID, variant_feature_name string) (int64, error)
}
type productRepository struct {
	db              *mongo.Database
	collection_name string
}

// AddProductVariant implements ProductRepository.
func (p *productRepository) AddProductVariant(ctx context.Context, productId primitive.ObjectID, variant domain.ProductVariant) (int64, error) {
	c, cancel := context.WithTimeout(ctx, time.Second*10)
	collection := p.db.Collection(p.collection_name)
	defer cancel()
	update := bson.M{
		"$push": bson.M{
			"product_variant": variant,
		},
	}

	// Perform the update using the product's _id
	res, err := collection.UpdateByID(c, productId, update)
	if err != nil {
		log.Error(err)
		fmt.Print(err)
		return 0, err
	}
	return res.MatchedCount, nil

}

// DeleteProductVariant implements ProductRepository.
func (p *productRepository) DeleteProductVariant(ctx context.Context, productId primitive.ObjectID, variant_feature_name string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	collection := p.db.Collection(p.collection_name)
	deleteItem := bson.M{
		"$pull": bson.M{
			"product_variant": bson.M{
				"variant_features": variant_feature_name, // Match variant based on the feature name
			},
		},
	}

	res, err := collection.UpdateByID(ctx, productId, deleteItem)
	if err != nil {
		log.Error(err)
		fmt.Print("error at line 80")
		fmt.Print(err)
		return 0, err
	}
	return res.ModifiedCount, nil
}

// UpdateProductVariant implements ProductRepository.
func (p *productRepository) UpdateProductVariant(ctx context.Context, productId primitive.ObjectID, updatedVariant domain.ProductVariant) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	collection := p.db.Collection(p.collection_name)

	// Filter to find the product by its ID and the variant by its feature name
	filter := bson.M{
		"_id":                              productId,
		"product_variant.variant_features": updatedVariant.VarientFeatures, // Match specific variant by feature name
	}

	// Update the matched variant with the new data
	update := bson.M{
		"$set": bson.M{
			"product_variant.$": updatedVariant, // The positional operator "$" updates the matched variant
		},
	}

	// Perform the update operation
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	// Return the number of modified documents
	return res.ModifiedCount, nil
}

// CreateProduct implements ProductRepository.
func (p *productRepository) CreateProduct(ctx context.Context, product *domain.ProductRequest) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	collection := p.db.Collection(p.collection_name)
	defer cancel()

	res, err := collection.InsertOne(ctx, bson.M{
		"description":     product.Description,
		"product_name":    product.ProductName,
		"image_link":      product.ImageLink,
		"product_variant": product.Variant,
	})
	if err != nil {
		return nil, err
	}

	prod := domain.Product{
		ProductId:   res.InsertedID.(primitive.ObjectID),
		ProductName: product.ProductName,
		Description: product.Description,
		ImageLink:   product.ImageLink,
		Variant:     product.Variant,
	}
	return &prod, nil
}

// DeleteProductById implements ProductRepository.
func (p *productRepository) DeleteProductById(ctx context.Context, id primitive.ObjectID) (rowAffect int64, err error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	collection := p.db.Collection(p.collection_name)
	defer cancel()

	res, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return res.DeletedCount, nil
}

// GetAllProducts implements ProductRepository.
func (p *productRepository) GetAllProducts(ctx context.Context) ([]*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	collection := p.db.Collection(p.collection_name)

	defer cancel()
	curprod, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer curprod.Close(ctx) // Ensure the cursor is closed after usage
	var prodlist []*domain.Product
	for curprod.Next(ctx) {
		var product domain.Product
		err := curprod.Decode(&product)
		if err != nil {
			return nil, err
		}

		prodlist = append(prodlist, &product)

	}
	// Check if the cursor encountered any errors while iterating
	if err := curprod.Err(); err != nil {
		return nil, err
	}

	return prodlist, nil
}

// GetProductById implements ProductRepository.
func (p *productRepository) GetProductById(ctx context.Context, id primitive.ObjectID) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	collection := p.db.Collection(p.collection_name)
	defer cancel()

	var prod domain.Product
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&prod)
	// if err != nil {
	// 	log.Error(err)
	// 	return nil, err
	// }
	if err != nil {
		// Check if the error is because the product wasn't found
		if errors.Is(err, mongo.ErrNoDocuments) {
			//log.Errorf("Product with ID %s not found", id.Hex())
			return nil, fmt.Errorf("product with ID %s not found", id.Hex())
		}
		// Log other errors that may have occurred
		//log.Errorf("Error retrieving product with ID %s: %v", id.Hex(), err)
		return nil, fmt.Errorf("error retrieving product with ID %s: %w", id.Hex(), err)
	}
	prod.ProductId = id

	return &prod, nil
}

// UpdateProductById implements ProductRepository.
func (p *productRepository) UpdateProductById(ctx context.Context, id primitive.ObjectID, updatedProduct *domain.ProductRequest) (int64, error) {
	c, cancel := context.WithTimeout(ctx, time.Second*5)
	collection := p.db.Collection(p.collection_name)
	defer cancel()
	res, err := collection.UpdateByID(c, id, bson.M{
		"$set": bson.M{
			"product_name":    updatedProduct.ProductName,
			"description":     updatedProduct.Description,
			"image_link":      updatedProduct.ImageLink,
			"product_variant": updatedProduct.Variant,
		},
	})
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return res.MatchedCount, nil
}

func NewProductRepository(db *mongo.Database, collection_name string) ProductRepository {
	return &productRepository{
		db:              db,
		collection_name: collection_name,
	}
}
