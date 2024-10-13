package appctx

import "go.mongodb.org/mongo-driver/mongo"

type Appcontext interface {
	GetConnectionToDatabase() *mongo.Database
	GetSecretKeyString() string
}

type appcontext struct {
	db        *mongo.Database
	secretkey string
}

func NewAppContext(db *mongo.Database, secreykey string) *appcontext {
	return &appcontext{
		db:        db,
		secretkey: secreykey,
	}
}
func (appctx *appcontext) GetConnectionToDatabase() *mongo.Database {
	return appctx.db
}

func (appctx *appcontext) GetSecretKeyString() string {
	return appctx.secretkey
}
