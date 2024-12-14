package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/domain"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"github.com/ydb-platform/ydb-go-sdk/v3/log"
)

type OrderHandler struct {
	OrderUsecase domain.OrderUsecase
	Env          *appconfig.Env
}

func (o *OrderHandler) CreatOrderHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("user_id")
		fmt.Printf("User id is that %v", userID)
		if !exists {
			fmt.Print("User ID not found in context")
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid user param",
				Error:   errors.New("Can not find user id in context").Error(),
			})
			return
		}
		var order_req domain.OrderReq
		if err := ctx.ShouldBindJSON(&order_req); err != nil {
			fmt.Print("User ID not found in context")
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Can not parsing the request",
				Error:   err.Error(),
			})
			return
		}
		hex_user_id, err := primitive.ObjectIDFromHex(userID.(string))
		if err != nil {
			fmt.Print("User ID not found in context")
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error converting user id string to object id",
				Error:   err.Error(),
			})
			return
		}

		fmt.Println(order_req.OrderItems)
		fmt.Println("#######################")
		//var order domain.Order
		order := domain.Order{
			ID:              primitive.NewObjectID(),
			CustomerID:      hex_user_id,
			CustomerName:    order_req.CustomerName,
			ShippingAddress: order_req.ShippingAddress,
			PhoneNumber:     order_req.PhoneNumber,
			Notes:           order_req.Notes,
			OrderItems:      order_req.OrderItems,
			TotalPrice:      0,
			PaymentInfo: domain.Payment{
				PaymentMethod: "cash",
				IsPaid:        0,
			},
			Email:          order_req.Email,
			ServiceType:    order_req.ServiceType,
			OrderStatus:    "pending",
			ShippingStatus: "pending",
			CreatedAt:      time.Now(),
		}
		// order.CustomerID = hex_user_id
		// order.ID = primitive.NewObjectID()
		// order.CreatedAt = time.Now()
		// order.UpdatedAt = time.Now()
		// order.OrderItems = order_req.OrderItems
		// order.PaymentInfo = order_req.PaymentInfo
		// order.ShippingAddress = order_req.ShippingAdress
		// order.OrderStatus = "pending"

		for _, row := range order_req.OrderItems {
			order.TotalPrice += row.Price * float64(row.BuyQuantity)
		}

		_, err = o.OrderUsecase.CreateOrder(ctx, order)
		if err != nil {
			fmt.Print("User ID not found in context")
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error happened while inserting to database",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusCreated, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusCreated,
				Message: "Successfully created order",
			},
			Data: order,
		})
	}
}
func (o *OrderHandler) GetAllOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orderlist, err := o.OrderUsecase.GetAllOrders(ctx)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error happened while getting order from database",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Get order list successfully",
			},
			Data: orderlist,
		})
	}
}
func (o *OrderHandler) UpdateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var updated_order domain.Order
		if err := ctx.ShouldBindJSON(&updated_order); err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Can not parsing body request",
				Error:   err.Error(),
			})
			return
		}
		updated_order.UpdatedAt = time.Now()
		res, err := o.OrderUsecase.UpdateOrder(ctx, updated_order)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while updade record in database",
				Error:   err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully update order",
			},
			Data: res,
		})

	}
}

func (o *OrderHandler) UpdateOrderStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var req domain.PaymentReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error parsing the json request",
				Error:   err.Error(),
			})
			return
		}
		reqhexid, err := primitive.ObjectIDFromHex(req.Order_id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Fail to convert id from objectid",
				Error:   err.Error(),
			})
			return
		}
		row_affected, err := o.OrderUsecase.UpdateOrderPaymentStatus(ctx, reqhexid, req.PaymentInfo.IsPaid)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error happened while querying database",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully update order payment status",
			},
			Data: gin.H{"row_affected": row_affected},
		})
	}
}
func (o *OrderHandler) GetOrderByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		order_param := ctx.Param("order_id")

		objID, err := primitive.ObjectIDFromHex(order_param)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error converting order id to object id",
				Error:   err.Error(),
			})
			return
		}
		order, err := o.OrderUsecase.GetOrderByID(ctx, objID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error happened while querying database",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully get order by id",
			},
			Data: order,
		})
	}
}
func (o *OrderHandler) GetOrderByCustomerID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		order_param := ctx.Query("customer_id")

		objID, err := primitive.ObjectIDFromHex(order_param)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error converting order id to object id",
				Error:   err.Error(),
			})
			return
		}
		order, err := o.OrderUsecase.GetOrdersByCustomerID(ctx, objID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error happened while querying database",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully get order by customer id",
			},
			Data: order,
		})
	}
}
