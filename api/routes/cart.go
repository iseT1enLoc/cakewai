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

func NewCartRoute(Env *appconfig.Env, timout time.Duration, db *mongo.Database, r *gin.RouterGroup) {
	cart_repo := repository.NewCartRepository(db, "cart")
	cart_handler := handlers.CartHandler{
		CartUseCase: usecase.NewCartUsecase(cart_repo, timout),
		Env:         Env,
	}
	r.GET("/cart:user_id", cart_handler.CreateCartByUserId())
	r.GET("/cart/items", cart_handler.GetAllItemsInCartByUserID())
	r.GET("/cart/", cart_handler.GetCartByUserID())
	r.DELETE("/cart/item?category={cart_id}&product_id={prod_id}&variant={variant}", cart_handler.RemoveItemFromCart())
	r.PUT("/cart/additem", cart_handler.AddCartItemIntoCart())
	r.PUT("/variant/:product_id", cart_handler.UpdateCartItemByID())
}
