package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/component/response"
	"cakewai/cakewai.com/domain"
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"github.com/ydb-platform/ydb-go-sdk/v3/log"
)

type CartHandler struct {
	CartUseCase domain.CartUsecase
	Env         *appconfig.Env
}

func (ch *CartHandler) CreateCartByUserId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id := ctx.Param("id")
		objId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid user param",
				Error:   err.Error(),
			})
		}
		err = ch.CartUseCase.CreateCartByUserId(ctx, objId)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error happened while insert cart into database",
				Error:   err.Error(),
			})
		}
		ctx.JSON(http.StatusCreated, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusCreated,
				Message: "Successfully create empty cart",
			},
			Data: nil,
		})

	}
}
func (ch *CartHandler) GetAllItemsInCartByUserID() gin.HandlerFunc {
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
		objhex, err := primitive.ObjectIDFromHex(userID.(string))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid user param",
				Error:   err.Error(),
			})
		}
		all_items_in_carts, err := ch.CartUseCase.GetAllItemsInCartByUserID(ctx, objhex)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while querying database",
				Error:   err.Error(),
			})
		}
		ctx.JSON(http.StatusCreated, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusCreated,
				Message: "Successfully create empty cart",
			},
			Data: all_items_in_carts,
		})
	}
}
func (ch *CartHandler) GetCartByUserID() gin.HandlerFunc {
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
		objhex, err := primitive.ObjectIDFromHex(userID.(string))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "invalid user param",
				Error:   err.Error(),
			})
		}
		cart, err := ch.CartUseCase.GetCartByUserID(ctx, objhex)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "can not convert into hexa id",
				Error:   err.Error(),
			})
		}
		ctx.JSON(http.StatusCreated, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusCreated,
				Message: "Successfully create empty cart",
			},
			Data: cart,
		})
	}
}

// http://localhost:8080/api/v1/items?category=books&price_min=10&price_max=50
func (ch *CartHandler) RemoveItemFromCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cart_id := ctx.Query("cart_id")
		productID := ctx.Query("product_id")
		variant := ctx.Query("variant")
		if productID == "" || variant == "" || cart_id == "" {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "cart id, productid or variant is empty",
				Error:   errors.New("productID and variant is empty").Error(),
			})
			return
		}
		cartid_hex, err1 := primitive.ObjectIDFromHex(cart_id)
		productid_hex, err2 := primitive.ObjectIDFromHex(productID)
		if err1 != nil || err2 != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "CartID or product id is invalid",
				Error:   err1.Error() + err2.Error(),
			})
			return
		}
		err := ch.CartUseCase.RemoveItemFromCart(ctx, cartid_hex, productid_hex, variant)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while removing data from cart",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusCreated, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusCreated,
				Message: "Successfully remove item from cart",
			},
			Data: nil,
		})
	}
}
func (ch *CartHandler) AddCartItemIntoCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var item domain.CartItem
		if err := ctx.ShouldBindJSON(&item); err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Can not parsing item",
				Error:   err.Error(),
			})
			return
		}
		item_id, err := ch.CartUseCase.AddCartItemIntoCart(ctx, item)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error happened while adding item into cart",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusCreated, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusCreated,
				Message: "Successfully create empty cart",
			},
			Data: item_id,
		})
	}
}
func (ch *CartHandler) UpdateCartItemByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var updated_item domain.CartItem
		if err := ctx.ShouldBindJSON(&updated_item); err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Can not parsing item",
				Error:   err.Error(),
			})
			return
		}
		updatedItem, err := ch.CartUseCase.UpdateCartItemByID(ctx, updated_item)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error happened while updating item into cart",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusCreated, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusCreated,
				Message: "Successfully update item in cart",
			},
			Data: updatedItem,
		})
	}
}
