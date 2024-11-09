package usecase

import (
	"cakewai/cakewai.com/domain"
	"cakewai/cakewai.com/repository"
	"context"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"
)

type productTypeUsecase struct {
	productTypeRepository repository.ProductTypeRepository
	timeout               time.Duration
}

// CreateProductType implements domain.ProductTypeUsecase.
func (p *productTypeUsecase) CreateProductType(context context.Context, pro_type domain.ProductType) error {
	err := p.productTypeRepository.CreateProductType(context, pro_type)
	log.Error(err)
	return err
}

// GetAllProductType implements domain.ProductTypeUsecase.
func (p *productTypeUsecase) GetAllProductType(context context.Context) ([]domain.ProductType, error) {
	panic("unimplemented")
}

// GetProductTypeById implements domain.ProductTypeUsecase.
func (p *productTypeUsecase) GetProductTypeById(context context.Context, product_type_id int) (*domain.ProductType, error) {
	panic("unimplemented")
}

// RemoveProductType implements domain.ProductTypeUsecase.
func (p *productTypeUsecase) RemoveProductType(context context.Context, product_type_id int) error {
	panic("unimplemented")
}

// UpdateProductType implements domain.ProductTypeUsecase.
func (p *productTypeUsecase) UpdateProductType(context context.Context, updated_product_type domain.ProductType) (*domain.ProductType, error) {
	panic("unimplemented")
}

func NewProductTypeUsecase(repo repository.ProductTypeRepository, timeout time.Duration) domain.ProductTypeUsecase {
	return &productTypeUsecase{
		productTypeRepository: repo,
		timeout:               timeout,
	}
}
