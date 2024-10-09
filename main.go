package main

import (
	"log"
	"net/http"
	"time"

	"cakewai/cakewai.com/api/middlewares"
	"cakewai/cakewai.com/api/routes"
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

	routes.SetUp(appcfg, time.Duration(appcfg.REFRESH_TOK_EXP), db, r)
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"data": "hello everyone, welcome to my chanel"})
	})
	//r.Run("localhost:8080")
	r.Run()
}
