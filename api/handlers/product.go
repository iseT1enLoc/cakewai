package handlers

import (
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/component/response"
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

			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error happened while parsing product json",
				Error:   err.Error(),
			})
			return
		}

		prod, err := pc.ProductUsecase.CreateProduct(ctx, &product)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error happened while inserting to database",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.SuccessResponse{
			Success: response.Success{
				ResponseFormat: response.ResponseFormat{
					Code:    http.StatusOK,
					Message: "Successfully create product",
				},
				Data: prod,
			},
		})
	}
}

func (pc *ProductHandler) GetProductById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		prodIDParam := ctx.Param("product_id")
		prodID, err := primitive.ObjectIDFromHex(prodIDParam)

		if err != nil {
			log.Error(errors.Errorf("Invalid product id %s , error %s", prodID, err))
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    200,
				Message: "Invalid product id",
				Error:   err.Error(),
			})
			return
		}

		prod, err := pc.ProductUsecase.GetProductById(ctx, prodID)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error happened while querying database..",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.SuccessResponse{
			Success: response.Success{
				ResponseFormat: response.ResponseFormat{
					Code:    http.StatusOK,
					Message: "",
				},
				Data: prod,
			},
		})
	}
}

func (pc *ProductHandler) UpdateProductById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		prodIDParam := ctx.Param("product_id")
		prodID, err := primitive.ObjectIDFromHex(prodIDParam)
		if err != nil {
			log.Error(errors.Errorf("Invalid product id %s , error %s", prodID, err))
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid product id",
				Error:   err.Error(),
			})
			return
		}
		var updatedProd *domain.ProductRequest
		if err := ctx.ShouldBindJSON(&updatedProd); err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while parsing json..",
				Error:   err.Error(),
			})
			return
		}

		rowAffected, err := pc.ProductUsecase.UpdateProductById(ctx.Request.Context(), prodID, updatedProd)
		if err != nil {
			fmt.Errorf("Error while update product by id")
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while update product by id",
				Error:   err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, response.SuccessResponse{
			Success: response.Success{
				ResponseFormat: response.ResponseFormat{
					Code:    http.StatusOK,
					Message: "Successfully update product by id",
				},
				Data: rowAffected,
			},
		})
	}
}
func (pc *ProductHandler) DeleteProductById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		prodParam := ctx.Param("product_id")
		prodID, err := primitive.ObjectIDFromHex(prodParam)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid product id",
				Error:   err.Error(),
			})
			return
		}
		err = pc.ProductUsecase.DeleteProductById(ctx, prodID)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while deleting product...",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.SuccessResponse{
			Success: response.Success{
				ResponseFormat: response.ResponseFormat{
					Code:    http.StatusOK,
					Message: "Successfully deleting product",
				},
				Data: nil,
			},
		})
	}
}
func (pc *ProductHandler) GetAllProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productlist, err := pc.ProductUsecase.GetAllProducts(ctx)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Error while querying database..",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.SuccessResponse{
			Success: response.Success{
				ResponseFormat: response.ResponseFormat{
					Code:    http.StatusOK,
					Message: "Successfully getting product list",
				},
				Data: productlist,
			},
		})
	}
}
func (pc *ProductHandler) AddProductVariant() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		prodIDparam := ctx.Param("product_id")
		product_id, err := primitive.ObjectIDFromHex(prodIDparam)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid product id..",
				Error:   err.Error(),
			})
			return
		}
		var prod_variant domain.ProductVariant
		if err := ctx.ShouldBindJSON(&prod_variant); err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Can not parsing product id",
				Error:   err.Error(),
			})
			return
		}
		_, err = pc.ProductUsecase.AddProductVariant(ctx, product_id, prod_variant)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error happened while adding product variant into db",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.SuccessResponse{
			Success: response.Success{
				ResponseFormat: response.ResponseFormat{
					Code:    http.StatusOK,
					Message: "Successfully getting adding product variant...",
				},
				Data: nil,
			},
		})
	}
}
func (pc *ProductHandler) DeleteProductVariant() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		prodParam := ctx.Param("product_id")
		id, err := primitive.ObjectIDFromHex(prodParam)
		var variant struct {
			VariantName string `json:"variant_name"`
		}
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid product variant id...",
				Error:   err.Error(),
			})
			return
		}
		if err := ctx.ShouldBindJSON(&variant); err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while parsing product variant name...",
				Error:   err.Error(),
			})
			return
		}
		_, err = pc.ProductUsecase.DeleteProductVariant(ctx, id, variant.VariantName)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while deleting product variant of this product id...",
				Error:   err.Error(),
			})
			return
		}
		prod, err := pc.ProductUsecase.GetProductById(ctx, id)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusInternalServerError,
				Message: "Error happened while getting the product by id..",
				Error:   "",
			})
			return
		}
		ctx.JSON(http.StatusOK, response.SuccessResponse{
			Success: response.Success{
				ResponseFormat: response.ResponseFormat{
					Code:    http.StatusOK,
					Message: "Successfully deleting product variant by name and here is the current product status...",
				},
				Data: prod,
			},
		})
	}
}
func (pc *ProductHandler) UpdateProductVarientByName() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		prodParam := ctx.Param("product_id")
		prodID, err := primitive.ObjectIDFromHex(prodParam)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while converting product id from hex to string...",
				Error:   err.Error(),
			})
			return
		}
		var updatedVariant domain.ProductVariant
		if err := ctx.ShouldBindJSON(&updatedVariant); err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while parsing product variant",
				Error:   err.Error(),
			})
		}
		_, err = pc.ProductUsecase.UpdateProductVariant(ctx, prodID, updatedVariant)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while updating product variant...",
				Error:   err.Error(),
			})
			return
		}
		prod, err := pc.ProductUsecase.GetProductById(ctx, prodID)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusInternalServerError, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.SuccessResponse{
			Success: response.Success{
				ResponseFormat: response.ResponseFormat{
					Code:    http.StatusOK,
					Message: "Successfully update product...",
				},
				Data: prod,
			},
		})

	}
}

func (pc *ProductHandler) GetProductByProductTypeID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("type_id")

		products, err := pc.ProductUsecase.GetProductByProductTypeID(ctx, id)
		if err != nil {
			log.Error(err)
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error happened while querying database, maybe invalid product type",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "Successfully get product by product type id",
			},
			Data: products,
		})
	}
}

func (pc *ProductHandler) SearchProducts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		query := ctx.Query("query")
		typeID := ctx.Query("type_id")           // Filter by ProductTypeID
		variantName := ctx.Query("variant_name") // Filter by Variant Name

		products, err := pc.ProductUsecase.SearchProducts(ctx, query, typeID, variantName)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.FailedResponse{
				Code:    http.StatusBadRequest,
				Message: "Error while getting desired products",
				Error:   err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    http.StatusOK,
				Message: "List of results",
			},
			Data: products,
		})
	}
}
func (pc *ProductHandler) FilterProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()

		// Extract query parameters
		sortField := c.Query("field") // Field to sort by (e.g., "variant.price" or "product_name")
		sortOrder := c.Query("order") // Sort order ("asc" or "desc")

		// Validate query parameters
		if sortField == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "field parameter is required"})
			return
		}
		if sortOrder != "asc" && sortOrder != "desc" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order parameter. Use 'asc' or 'desc'."})
			return
		}

		// Fetch sorted products from the usecase
		products, err := pc.ProductUsecase.FetchSortedProducts(ctx, sortField, sortOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Respond with the sorted products
		c.JSON(http.StatusOK, response.Success{
			ResponseFormat: response.ResponseFormat{
				Code:    200,
				Message: "Successfully get sorted products list",
			},
			Data: products,
		})
	}
}
