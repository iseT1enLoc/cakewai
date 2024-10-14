package usecase

import (
	"cakewai/cakewai.com/domain"
	"cakewai/cakewai.com/repository"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type productUsecase struct {
	repository repository.ProductRepository
	timeout    time.Duration
}

// AddProductVariant implements domain.ProductRepository.
func (p productUsecase) AddProductVariant(ctx context.Context, productId primitive.ObjectID, variant domain.ProductVariant) (*domain.Product, error) {
	panic("unimplemented")
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
func (p productUsecase) DeleteProductVariant(ctx context.Context, productId primitive.ObjectID, variantId primitive.ObjectID) (*domain.Product, error) {
	panic("unimplemented")
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
func (p productUsecase) UpdateProductVariant(ctx context.Context, productId primitive.ObjectID, variantId primitive.ObjectID, updatedVariant domain.ProductVariant) (*domain.Product, error) {
	panic("unimplemented")
}

func NewProductUsecase(repo repository.ProductRepository, timeout time.Duration) domain.ProductUsecase {
	return &productUsecase{
		repository: repo,
		timeout:    timeout,
	}
}
