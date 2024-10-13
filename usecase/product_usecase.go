package usecase

import (
	"cakewai/cakewai.com/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type productUsecase struct {
	db              mongo.Database
	collection_name string
}

func NewProductUsecase(database mongo.Database, collection string) domain.ProductRepository {

}
