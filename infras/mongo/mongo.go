package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	appconfig "cakewai/cakewai.com/component/appcfg"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectWithMongodb(appcfg *appconfig.Env) (*mongo.Client, error) {
	timeTry := time.Second * 20 //time to connect to database
	connectingToMongoDB := func(appcfg *appconfig.Env) (*mongo.Client, error) {
		mongo_db_url := appcfg.DATABASE_URL

		// load .env file
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
		// set mongodb connection string
		if mongo_db_url == "" {
			log.Fatal("MONGODB_URI is not set")
		}
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongo_db_url))
		if err != nil {
			log.Fatal(err)
		}

		// check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("Connected to mongodb")
		}

		return client, nil
	}
	print("Line 44 connect databasez")
	deadline := time.Now().Add(time.Duration(timeTry))
	var dbclient *mongo.Client
	var err error

	for time.Now().Before(deadline) {
		fmt.Print("line 50")
		log.Println("CONNECT to database.....")
		dbclient, err = connectingToMongoDB(appcfg)
		if err == nil {
			//fmt.Printf("Database name: %v", db.Name())
			fmt.Print("line 55 conect database")
			return dbclient, nil
		}
		time.Sleep(time.Second)
	}
	fmt.Print("Line 56")
	return nil, fmt.Errorf("Error while connecting to database...[error]: %v", err)
}

// OpenCollection get collection
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cakewai").Collection(collectionName)
	return collection
}
