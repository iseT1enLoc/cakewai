package usecase

import (
	"cakewai/cakewai.com/domain"
	"cakewai/cakewai.com/repository"
	"context"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type productUsecase struct {
	repository repository.ProductRepository
	timeout    time.Duration
}

// CreateProductType implements domain.ProductTypeUsecase.
func (p *productUsecase) CreateProductType(context context.Context) error {
	panic("unimplemented")
}

// GetAllProductType implements domain.ProductTypeUsecase.
func (p *productUsecase) GetAllProductType(context context.Context) ([]domain.ProductType, error) {
	panic("unimplemented")
}

// GetProductTypeById implements domain.ProductTypeUsecase.
func (p *productUsecase) GetProductTypeById(context context.Context, product_type_id int) (*domain.ProductType, error) {
	panic("unimplemented")
}

// RemoveProductType implements domain.ProductTypeUsecase.
func (p *productUsecase) RemoveProductType(context context.Context, product_type_id int) error {
	panic("unimplemented")
}

// UpdateProductType implements domain.ProductTypeUsecase.
func (p *productUsecase) UpdateProductType(context context.Context, updated_product_type domain.ProductType) (*domain.ProductType, error) {
	panic("unimplemented")
}

// AddProductVariant implements domain.ProductRepository.
func (p productUsecase) AddProductVariant(ctx context.Context, productId primitive.ObjectID, variant domain.ProductVariant) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	rowAffected, err := p.repository.AddProductVariant(ctx, productId, variant)
	if err != nil {
		log.Error(err)

		return 0, err
	}
	return rowAffected, nil
}

// CreateProduct implements domain.ProductRepository.
func (p productUsecase) CreateProduct(ctx context.Context, product *domain.ProductRequest) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	return p.repository.CreateProduct(ctx, product)
}

// DeleteProductById implements domain.ProductRepository.
func (p productUsecase) DeleteProductById(ctx context.Context, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	_, err := p.repository.DeleteProductById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteProductVariant implements domain.ProductRepository.
func (p productUsecase) DeleteProductVariant(ctx context.Context, productId primitive.ObjectID, variant string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	_, err := p.repository.DeleteProductVariant(ctx, productId, variant)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return 1, nil
}

// GetAllProducts implements domain.ProductRepository.
func (p productUsecase) GetAllProducts(ctx context.Context) ([]*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	prodlist, err := p.repository.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}
	return prodlist, nil
}

// GetProductById implements domain.ProductRepository.
func (p productUsecase) GetProductById(ctx context.Context, id primitive.ObjectID) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	return p.repository.GetProductById(ctx, id)
}

// UpdateProductById implements domain.ProductRepository.
func (p productUsecase) UpdateProductById(ctx context.Context, id primitive.ObjectID, updatedProduct *domain.ProductRequest) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	return p.repository.UpdateProductById(ctx, id, updatedProduct)
}

// UpdateProductVariant implements domain.ProductRepository.
func (p productUsecase) UpdateProductVariant(ctx context.Context, productId primitive.ObjectID, updatedVariant domain.ProductVariant) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	return p.repository.UpdateProductVariant(ctx, productId, updatedVariant)
}

func NewProductUsecase(repo repository.ProductRepository, timeout time.Duration) domain.ProductUsecase {
	return &productUsecase{
		repository: repo,
		timeout:    timeout,
	}
}
