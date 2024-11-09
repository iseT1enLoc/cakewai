package repository

import (
	apperror "cakewai/cakewai.com/component/apperr"
	"cakewai/cakewai.com/domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductTypeRepository interface {
	CreateProductType(context context.Context, product_type domain.ProductType) error
	GetAllProductType(context context.Context) ([]domain.ProductType, error)
	GetProductTypeById(context context.Context, product_type_id int) (*domain.ProductType, error)
	RemoveProductType(context context.Context, product_type_id int) error
	UpdateProductType(context context.Context, updated_product_type domain.ProductType) (*domain.ProductType, error)
}
type productTypeRepository struct {
	db              *mongo.Database
	collection_name string
}

// DONE
func (p *productTypeRepository) CreateProductType(context context.Context, product_type domain.ProductType) error {
	collection := p.db.Collection(p.collection_name)

	values := bson.M{"_id": primitive.NewObjectID(),
		"type_name":   product_type.Product_type_name,
		"description": product_type.Description}
	_, err := collection.InsertOne(context, values)
	return err
}

// GetAllProductType implements ProductTypeRepository.
func (p *productTypeRepository) GetAllProductType(context context.Context) ([]domain.ProductType, error) {
	return []domain.ProductType{}, apperror.ErrFailedGenerateJWT
}

// GetProductTypeById implements ProductTypeRepository.
func (p *productTypeRepository) GetProductTypeById(context context.Context, product_type_id int) (*domain.ProductType, error) {
	return &domain.ProductType{}, apperror.ErrCodeExchangeWrong
}

// RemoveProductType implements ProductTypeRepository.
func (p *productTypeRepository) RemoveProductType(context context.Context, product_type_id int) error {
	return apperror.ErrCodeExchangeWrong
}

// UpdateProductType implements ProductTypeRepository.
func (p *productTypeRepository) UpdateProductType(context context.Context, updated_product_type domain.ProductType) (*domain.ProductType, error) {
	return &domain.ProductType{}, apperror.ErrCodeExchangeWrong
}

func NewProductTypeRepository(db *mongo.Database, collection_name string) ProductTypeRepository {
	return &productTypeRepository{
		db:              db,
		collection_name: collection_name,
	}
}
