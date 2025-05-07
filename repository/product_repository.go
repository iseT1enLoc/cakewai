package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cakewai/cakewai.com/domain"
	"cakewai/cakewai.com/internals/utils"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	// Create a new product: DONE-------
	CreateProduct(ctx context.Context, product *domain.ProductRequest) (*domain.Product, error)

	// Get product by ID:DONE------------
	GetProductById(ctx context.Context, id primitive.ObjectID) (*domain.Product, error)

	// Get product by ID:DONE------------
	GetProductByProductTypeID(ctx context.Context, id string) ([]*domain.Product, error)

	// Get all products: DONE------------
	GetAllProducts(ctx context.Context) ([]*domain.Product, error)

	// Update product by ID:DONE----------
	UpdateProductById(ctx context.Context, id primitive.ObjectID, updatedProduct *domain.ProductRequest) (int64, error)

	// Delete product by ID: DONE-------
	DeleteProductById(ctx context.Context, id primitive.ObjectID) (rowAffect int64, err error)

	// Add a variant to a product:DONE----------
	AddProductVariant(ctx context.Context, productId primitive.ObjectID, variant domain.ProductVariant) (int64, error)

	// Update a product variant
	UpdateProductVariant(ctx context.Context, productId primitive.ObjectID, updatedVariant domain.ProductVariant) (int64, error)

	// Delete a product variant
	DeleteProductVariant(ctx context.Context, productId primitive.ObjectID, variant_feature_name string) (int64, error)

	//Search product
	SearchProduct(ctx context.Context, query string, type_id string, variant string) ([]*domain.Product, error)

	//Fetch sorted products
	FetchSortedProducts(ctx context.Context, sortField string, sortOrder string) ([]*domain.Product, error)
}
type productRepository struct {
	db              *mongo.Database
	collection_name string
}

// FetchSortedProducts implements ProductRepository.
func (p *productRepository) FetchSortedProducts(ctx context.Context, sortField string, sortOrder string) ([]*domain.Product, error) {
	collection := p.db.Collection("products")

	// Determine sort order
	order := 1
	if sortOrder == "desc" {
		order = -1
	}

	var pipeline []bson.M

	// Check if sorting by variant field, specifically the first variant's price
	if sortField == "variant.price" {
		pipeline = []bson.M{
			// Add a new field with the price of the first variant (index 0) for sorting
			{"$addFields": bson.M{
				"first_variant_price": bson.M{
					"$arrayElemAt": []interface{}{"$product_variant.price", 0},
				},
			}},
			// Sort by the new field (first variant's price)
			{"$sort": bson.M{"first_variant_price": order}},
		}
	} else {
		// Default sort by the main product fields
		pipeline = []bson.M{
			{"$sort": bson.M{sortField: order}},
		}
	}

	// Perform aggregation
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate products: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode results into a slice of Product
	var products []*domain.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, nil
}

// SearchProduct implements ProductRepository.
func (p *productRepository) SearchProduct(ctx context.Context, query string, type_id string, variant string) ([]*domain.Product, error) {
	collection := p.db.Collection(p.collection_name)

	// Build MongoDB filter
	filter := bson.M{}

	// General text search on product name and description
	if query != "" {
		filter["$or"] = []bson.M{
			{"product_name": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
		}
	}

	// Filter by product type ID
	if type_id != "" {
		filter["product_type_id"] = type_id
	}

	// Filter by variant name
	if variant != "" {
		filter["product_variant.variant_features"] = bson.M{"$regex": variant, "$options": "i"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Perform the query
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		err := errors.New("Error searching products")
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var products []*domain.Product
	if err = cursor.All(ctx, &products); err != nil {
		err := errors.New("Error decoding search results")

		return nil, err
	}
	return products, nil

}

// GetProductByProductTypeID implements ProductRepository.
func (p *productRepository) GetProductByProductTypeID(ctx context.Context, id string) ([]*domain.Product, error) {
	collection := p.db.Collection(p.collection_name)
	var products []*domain.Product
	cur, err := collection.Find(ctx, bson.M{"product_type_id": id})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	fmt.Println("LINE 60 GET PRODUCT BY PRODUCT TYPE ID")

	err = cur.All(ctx, &products)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return products, nil
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
		"product_type_id": product.ProductTypeID,
	})
	if err != nil {
		return nil, err
	}

	prod := domain.Product{
		ProductId:     res.InsertedID.(primitive.ObjectID),
		ProductName:   product.ProductName,
		Description:   product.Description,
		ImageLink:     product.ImageLink,
		ProductTypeID: product.ProductTypeID,
		Variant:       product.Variant,
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
	//ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	collection := p.db.Collection(p.collection_name)

	//defer cancel()
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
		product.Slug = utils.Slugify(product.ProductName)
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
			"product_type_id": updatedProduct.ProductTypeID,
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
