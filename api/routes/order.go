package routes

import (
	"cakewai/cakewai.com/api/handlers"
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/domain"
	"cakewai/cakewai.com/repository"
	"cakewai/cakewai.com/usecase"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository interface {
	CreateOrder(context context.Context, order domain.Order) (*domain.Order, error)
	GetAllOrders(context context.Context) ([]*domain.Order, error)
	UpdateOrder(context context.Context, updatedOrder domain.Order) (*domain.Order, error)
	UpdateOrderStatus(context context.Context, order_id primitive.ObjectID, is_paid int) (int, error)
	GetOrderByID(context context.Context, ID primitive.ObjectID) (*domain.Order, error)
}

func NewOrderRoute(Env *appconfig.Env, timeout time.Duration, db *mongo.Database, r *gin.RouterGroup) {
	order_repo := repository.NewOrderRepository(db, "orders")
	order_handler := handlers.OrderHandler{
		OrderUsecase: usecase.NewOrderUsecase(order_repo, timeout),
		Env:          Env,
	}
	r.GET("/orders", order_handler.GetAllOrders())
	r.GET("/order/:id", order_handler.GetOrderByID())
	r.POST("/order", order_handler.CreatOrderHandler())
	r.PUT("/order/update", order_handler.UpdateOrder())
	r.PATCH("/order/paystatus", order_handler.UpdateOrderStatus())
}
