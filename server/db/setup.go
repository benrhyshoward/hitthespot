package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var mongoUserCollection *mongo.Collection

func init() {
	connectionString := os.Getenv("MONGO_CONNECTION_STRING")
	if connectionString == "" {
		log.Fatal("MONGO_CONNECTION_STRING not provided")
	}
	clientOptions := options.Client().ApplyURI(connectionString)

	var err error
	mongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	mongoUserCollection = mongoClient.Database("db").Collection("users")

	log.Print("Connected to MongoDB")
}
