package main

import (
	"log"
	"net/http"

	"cakewai/cakewai.com/api/middlewares"
	appconfig "cakewai/cakewai.com/component/appcfg"
	appctx "cakewai/cakewai.com/component/appcontext"
	"cakewai/cakewai.com/infras/postgres"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	appcfg, err := appconfig.LoadEnv()
	if err != nil {
		log.Fatalf("Error happened while loading config %v", err)
	}

	db, err := postgres.ConnectToDatabasein20s(appcfg)
	if err != nil {
		log.Fatalf("Error happened while connect to database %v", err)
	}
	appctx := appctx.NewAppContext(db, appcfg.SECRET_KEY)

	r.Use(middlewares.CORS())
	r.Use(middlewares.Recover(appctx))

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"data": "hello everyone, welcome to my chanel"})
	})
	r.Run("localhost:8080")
}
