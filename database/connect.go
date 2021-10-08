package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func Connect() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MongoDB uri didn't set. Please check .env file")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	DB = client.Database("sample_todos")
}
