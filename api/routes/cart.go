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
	//r.GET("/cart/:user_id", cart_handler.CreateCartByUserId())//DONE
	r.POST("/cart/additem", cart_handler.AddCartItemIntoCart()) //DONE

	r.GET("/cart/items", cart_handler.GetAllItemsInCartByUserID()) //DONE

	r.GET("/cart/", cart_handler.GetCartByUserID()) //DONE

	///api/carts/{cart_id}/items/{product_id}
	r.PUT("/cart/item", cart_handler.RemoveItemFromCart()) //DONE
	r.PUT("/cart/update_cart", cart_handler.UpdateCartItemByID())
}
