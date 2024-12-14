package repository

import (
	"cakewai/cakewai.com/domain"
	"context"
	"fmt"

	"github.com/ydb-platform/ydb-go-sdk/v3/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository interface {
	CreateOrder(context context.Context, order domain.Order) (*domain.Order, error)
	GetAllOrders(context context.Context) ([]*domain.Order, error)
	UpdateOrder(context context.Context, updatedOrder domain.Order) (*domain.Order, error)
	UpdateOrderStatus(context context.Context, order_id primitive.ObjectID, is_paid int) (int, error)
	GetOrderByID(context context.Context, ID primitive.ObjectID) (*domain.Order, error)
	GetOrdersByCustomerID(context context.Context, CustomerID primitive.ObjectID) ([]*domain.Order, error)
}

type orderRepository struct {
	db              *mongo.Database
	collection_name string
}

func (o *orderRepository) GetOrdersByCustomerID(ctx context.Context, CustomerID primitive.ObjectID) ([]*domain.Order, error) {
	collection := o.db.Collection(o.collection_name)

	// Filter to match the CustomerID
	filter := bson.M{"customer_id": CustomerID}

	// Find all matching documents
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		fmt.Printf("Failed to fetch orders by Customer ID: %v\n", err)
		return nil, fmt.Errorf("could not fetch orders: %w", err)
	}
	defer cursor.Close(ctx) // Ensure cursor is closed when done

	// Decode the documents into a slice of orders
	var orders []*domain.Order
	for cursor.Next(ctx) {
		var order domain.Order
		if err := cursor.Decode(&order); err != nil {
			fmt.Printf("Failed to decode order: %v\n", err)
			return nil, fmt.Errorf("could not decode order: %w", err)
		}
		// Append the pointer of the decoded order
		orders = append(orders, &order)
	}

	// Check if no orders were found
	if len(orders) == 0 {
		return nil, fmt.Errorf("no orders found for customer_id %v", CustomerID)
	}

	return orders, nil
}

// DONE
func (o *orderRepository) UpdateOrderStatus(context context.Context, order_id primitive.ObjectID, is_paid int) (int, error) {
	collection := o.db.Collection(o.collection_name)
	filter := bson.M{"_id": order_id}
	updated_payment_status := bson.M{"$set": bson.M{
		"payment_info.is_paid": is_paid,
	}}
	res, err := collection.UpdateOne(context, filter, updated_payment_status)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return int(res.MatchedCount), nil
}

// DONE
func (o *orderRepository) CreateOrder(context context.Context, order domain.Order) (*domain.Order, error) {
	collection := o.db.Collection(o.collection_name)
	_, err := collection.InsertOne(context, order)
	if err != nil {
		fmt.Printf("Failed to create order: %v", err)
		return nil, fmt.Errorf("could not insert order: %w", err)
	}
	return &order, nil
}

// DONE
func (o *orderRepository) GetAllOrders(context context.Context) ([]*domain.Order, error) {
	collection := o.db.Collection(o.collection_name)
	list_order_cursor, err := collection.Find(context, bson.D{})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer list_order_cursor.Close(context)

	var list_order []*domain.Order
	err = list_order_cursor.All(context, &list_order)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return list_order, nil
}

// DONE
func (o *orderRepository) GetOrderByID(context context.Context, ID primitive.ObjectID) (*domain.Order, error) {
	collection := o.db.Collection(o.collection_name)
	var expected_order domain.Order
	err := collection.FindOne(context, bson.M{"_id": ID}).Decode(&expected_order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("order not found: %w", err)
		}
		fmt.Printf("Failed to fetch order by ID: %v", err)
		return nil, fmt.Errorf("could not fetch order: %w", err)
	}
	return &expected_order, nil
}

// DONE
func (o *orderRepository) UpdateOrder(context context.Context, updatedOrder domain.Order) (*domain.Order, error) {
	collection := o.db.Collection(o.collection_name)
	updatedorder := bson.M{"$set": updatedOrder}
	res, err := collection.UpdateOne(context, bson.M{"_id": updatedOrder.ID}, updatedorder)
	if err != nil {
		fmt.Printf("Failed to update order: %v", err)
		return nil, fmt.Errorf("could not update order: %w", err)
	}
	if res.MatchedCount == 0 {
		return nil, fmt.Errorf("order not found")
	}
	return &updatedOrder, nil

}

func NewOrderRepository(db *mongo.Database, collection_name string) OrderRepository {
	return &orderRepository{
		db:              db,
		collection_name: collection_name,
	}
}
