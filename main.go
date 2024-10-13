package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cakewai/cakewai.com/api/middlewares"
	"cakewai/cakewai.com/api/routes"
	appconfig "cakewai/cakewai.com/component/appcfg"
	appctx "cakewai/cakewai.com/component/appcontext"
	mongodb "cakewai/cakewai.com/infras/mongo"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	appcfg, err := appconfig.LoadEnv()
	if err != nil {
		log.Fatalf("Error happened while loading config %v", err)
	}
	fmt.Print("Line 22")
	client, err := mongodb.ConnectWithMongodb(appcfg)
	if err != nil {
		log.Fatalf("Error happened while connect to database %v", err)
	}
	db := client.Database(os.Getenv("DB_NAME"))
	appctx := appctx.NewAppContext(db, appcfg.SECRET_KEY)

	r.Use(middlewares.CORS())
	r.Use(middlewares.Recover(appctx))

	routes.SetUp(appcfg, time.Duration(appcfg.REFRESH_TOK_EXP), db, r)
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"data": "hello everyone, welcome to my chanel"})
	})
	r.Run("localhost:8080")
	//r.Run()
}
