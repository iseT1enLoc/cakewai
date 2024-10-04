package appctx

import "database/sql"

type Appcontext interface {
	GetConnectionToDatabase() *sql.DB
	GetSecretKeyString() string
}

type appcontext struct {
	db        *sql.DB
	secretkey string
}

func NewAppContext(db *sql.DB, secreykey string) *appcontext {
	return &appcontext{
		db:        db,
		secretkey: secreykey,
	}
}
func (appctx *appcontext) GetConnectionToDatabase() *sql.DB {
	return appctx.db
}

func (appctx *appcontext) GetSecretKeyString() string {
	return appctx.secretkey
}
