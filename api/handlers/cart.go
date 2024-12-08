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
)

type CartHandler struct {
	CartUseCase domain.CartUsecase
	Env         *appconfig.Env
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
		fmt.Printf("Object id hexa from the cart get all items handler %v", objhex)
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
				Code:    http.StatusOK,
				Message: "Successfully get all cart items",
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
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusCreated,
				Message: "Successfully get cart by user id",
			},
			Data: cart,
		})
	}
}

// http://localhost:8080/api/v1/items?category=books&price_min=10&price_max=50
func (ch *CartHandler) RemoveItemFromCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Extract query parameters

		productID := ctx.Query("product_id")
		variant := ctx.Query("variant")

		// Validate required parameters
		if productID == "" || variant == "" {
			fmt.Printf("Missing parameters: product_id=%s, variant=%s", productID, variant)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Missing required parameters: product_id, or variant",
				Error:   "product_id, and variant must not be empty",
			})
			return
		}

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

		productIDHex, err := primitive.ObjectIDFromHex(productID)
		if err != nil {
			fmt.Printf("Invalid product_id: %s", productID)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid product_id",
				Error:   err.Error(),
			})
			return
		}

		// Call the use case to remove the item from the cart
		err = ch.CartUseCase.RemoveItemFromCart(ctx.Request.Context(), objhex, productIDHex, variant)
		if err != nil {
			fmt.Printf("Error removing item from cart: %v", err)
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to remove item from cart",
				Error:   err.Error(),
			})
			return
		}

		// Respond with success
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully removed item from cart",
			},
			Data: nil,
		})
	}
}

func (ch *CartHandler) AddCartItemIntoCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Bind the incoming JSON request body to the CartItem struct
		var item domain.CartItem
		if err := ctx.ShouldBindJSON(&item); err != nil {
			fmt.Printf("Error parsing item: %v", err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid input data, unable to parse cart item",
				Error:   err.Error(),
			})
			return
		}

		// Check if mandatory fields are provided
		if item.ProductId.IsZero() || item.Variant == "" || item.Price <= 0 || item.BuyQuantity <= 0 {
			fmt.Printf("Invalid cart item data: %+v", item)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Missing required fields in cart item",
				Error:   "ProductId, variant, price, and buyQuantity must be valid",
			})
			return
		}

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
		// Call the use case to add the cart item
		cart_item, err := ch.CartUseCase.AddCartItemIntoCart(ctx, objhex, item)
		if err != nil {
			fmt.Printf("Error adding item to cart: %v", err)
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Error while adding item into cart",
				Error:   err.Error(),
			})
			return
		}

		// Successful response
		ctx.JSON(http.StatusCreated, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully added item to cart",
			},
			Data: cart_item,
		})
	}
}

func (ch *CartHandler) UpdateCartItemByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Parse the updated cart item from the request body
		var updatedItem domain.CartItem
		if err := ctx.ShouldBindJSON(&updatedItem); err != nil {
			fmt.Printf("Error parsing updated item: %v", err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid input data, unable to parse cart item",
				Error:   err.Error(),
			})
			return
		}

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
		// Call the use case to update the cart item
		updatedItemResponse, err := ch.CartUseCase.UpdateAnCartItemByUserID(ctx, objhex, updatedItem)
		if err != nil {
			fmt.Printf("Error updating cart item: %v", err)
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Error updating item in the cart",
				Error:   err.Error(),
			})
			return
		}

		// Return the updated item in the response
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully updated item in cart",
			},
			Data: updatedItemResponse,
		})
	}
}
