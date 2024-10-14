package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/domain"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"github.com/go-faster/errors"
	"github.com/ydb-platform/ydb-go-sdk/v3/log"
)

type ProductHandler struct {
	ProductUsecase domain.ProductUsecase
	Env            *appconfig.Env
}

func (pc *ProductHandler) CreateProductHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var product domain.ProductRequest

		if err := ctx.ShouldBindJSON(&product); err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error happened while parsing product"})
			return
		}

		prod, err := pc.ProductUsecase.CreateProduct(ctx, &product)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error while creating product"})
		}
		ctx.JSON(http.StatusOK, gin.H{"data": prod})
	}
}

func (pc *ProductHandler) GetProductById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		prodIDParam := ctx.Param("product_id")
		prodID, err := primitive.ObjectIDFromHex(prodIDParam)

		if err != nil {
			log.Error(errors.Errorf("Invalid product id %s , error %s", prodID, err))
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID format"})
			return
		}

		prod, err := pc.ProductUsecase.GetProductById(ctx, prodID)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error happened while trying to get product by id from repositpository"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": prod})
	}
}

func (pc *ProductHandler) UpdateProductById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		prodIDParam := ctx.Param("product_id")
		prodID, err := primitive.ObjectIDFromHex(prodIDParam)
		if err != nil {
			log.Error(errors.Errorf("Invalid product id %s , error %s", prodID, err))
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID format"})
			return
		}
		var updatedProd *domain.ProductRequest
		if err := ctx.ShouldBindJSON(&updatedProd); err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not bind body with json..."})
			return
		}

		rowAffected, err := pc.ProductUsecase.UpdateProductById(ctx.Request.Context(), prodID, updatedProd)
		if err != nil {
			fmt.Errorf("Error while update product by id")
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		}

		ctx.JSON(http.StatusOK, rowAffected)
	}
}
func (pc *ProductHandler) DeleteProductById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		prodParam := ctx.Param("product_id")
		prodID, err := primitive.ObjectIDFromHex(prodParam)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Product id"})
			return
		}
		err = pc.ProductUsecase.DeleteProductById(ctx, prodID)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error while deleting product"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "successfully delete product"})
	}
}
func (pc *ProductHandler) GetAllProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productlist, err := pc.ProductUsecase.GetAllProducts(ctx)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error while get all product from database"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "get all products successfully", "data": productlist})
	}
}
