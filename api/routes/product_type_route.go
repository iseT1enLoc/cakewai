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

func NewProductTypeRoute(Env *appconfig.Env, timout time.Duration, db *mongo.Database, r *gin.RouterGroup) {
	prod_repo := repository.NewProductTypeRepository(db, "product_type")
	prod_handler := handlers.ProductTypeHandler{
		Product_type_usecase: usecase.NewProductTypeUsecase(prod_repo, timout),
		Env:                  Env,
	}
	r.POST("/product_type", prod_handler.CreateProductType())

}
