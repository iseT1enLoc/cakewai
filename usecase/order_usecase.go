package usecase

import (
	"cakewai/cakewai.com/domain"
	"cakewai/cakewai.com/repository"
	"context"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type orderUsecase struct {
	order_repository repository.OrderRepository
	timeout          time.Duration
}

// DeleteOrder implements domain.OrderUsecase.
func (o *orderUsecase) DeleteOrder(context context.Context, order_id primitive.ObjectID) error {
	err := o.order_repository.DeleteOrder(context, order_id)
	return err
}

// GetOrderByCustomerID implements domain.OrderUsecase.
func (o *orderUsecase) GetOrdersByCustomerID(context context.Context, CustomerID primitive.ObjectID) ([]*domain.Order, error) {
	order, err := o.order_repository.GetOrdersByCustomerID(context, CustomerID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return order, nil
}

// DONE
func (o *orderUsecase) CreateOrder(context context.Context, order domain.Order) (*domain.Order, error) {
	res, err := o.order_repository.CreateOrder(context, order)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return res, nil
}

// DONE
func (o *orderUsecase) GetAllOrders(context context.Context) ([]*domain.Order, error) {
	order_list, err := o.order_repository.GetAllOrders(context)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return order_list, nil
}

// DONE
func (o *orderUsecase) GetOrderByID(context context.Context, ID primitive.ObjectID) (*domain.Order, error) {
	order, err := o.order_repository.GetOrderByID(context, ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return order, nil
}

// DONE
func (o *orderUsecase) UpdateOrder(context context.Context, updatedOrder domain.Order) (*domain.Order, error) {
	res, err := o.order_repository.UpdateOrder(context, updatedOrder)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return res, nil
}

// DONE
func (o *orderUsecase) UpdateOrderPaymentStatus(context context.Context, order_id primitive.ObjectID, is_paid int) (int, error) {
	row_affected, err := o.order_repository.UpdateOrderStatus(context, order_id, is_paid)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return row_affected, nil
}

func NewOrderUsecase(order_repository repository.OrderRepository, timeout time.Duration) domain.OrderUsecase {
	return &orderUsecase{
		order_repository: order_repository,
		timeout:          timeout,
	}
}
