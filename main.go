package main

import (
	"cakewai/cakewai.com/api/middlewares"
	"cakewai/cakewai.com/api/routes"
	appconfig "cakewai/cakewai.com/component/appcfg"
	appctx "cakewai/cakewai.com/component/appcontext"
	mongodb "cakewai/cakewai.com/infras/mongo"
	"fmt"
	"log"
	"os"
	"time"

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
	print(err)
	if err != nil {
		log.Fatalf("Error happened while connect to database %v", err)
	}
	db := client.Database(os.Getenv("DB_NAME"))
	ctx := appctx.NewAppContext(db, appcfg.SECRET_KEY)

	r.Use(middlewares.CORS())
	r.Use(middlewares.Recover(ctx))

	routes.SetUp(appcfg, time.Second*3600, db, r)
	r.Use(middlewares.TraceMiddleware("root middleware"))
	// r.GET("/", func(ctx *gin.Context) {
	// 	ctx.JSON(http.StatusOK, gin.H{"data": "hello everyone, welcome to my chanel"})
	// })
	//r.Run("localhost:8080")
	r.Run()

}

//mongodb+srv://hoaibao123456:9DQ18WlqmeoF9f7u@cakewai.auh9m.mongodb.net/?retryWrites=true&w=majority&appName=cakewai
