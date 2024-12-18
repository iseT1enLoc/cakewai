package routes

import (
	"cakewai/cakewai.com/api/handlers"
	appconfig "cakewai/cakewai.com/component/appcfg"
	"cakewai/cakewai.com/repository"
	"cakewai/cakewai.com/usecase"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewProductRoute(Env *appconfig.Env, timout time.Duration, db *mongo.Database, r *gin.RouterGroup) {
	prod_repo := repository.NewProductRepository(db, "products")
	prod_handler := handlers.ProductHandler{
		ProductUsecase: usecase.NewProductUsecase(prod_repo, timout),
		Env:            Env,
	}
	r.POST("/product", prod_handler.CreateProductHandler())
	r.GET("/product/:product_id", prod_handler.GetProductById())
	r.GET("/products/:type_id", prod_handler.GetProductByProductTypeID())
	r.GET("/products/", prod_handler.GetAllProducts())
	r.PUT("/product/:product_id", prod_handler.UpdateProductById())
	r.POST("/variant/:product_id", prod_handler.AddProductVariant())
	r.PUT("/variant/:product_id", prod_handler.UpdateProductVarientByName())
	r.DELETE("/product/:product_id", prod_handler.DeleteProductById())
	r.DELETE("/variant/:product_id", prod_handler.DeleteProductVariant())
	r.GET("/product/search", prod_handler.SearchProducts())
	r.GET("/products/sort", prod_handler.FilterProduct())
}
