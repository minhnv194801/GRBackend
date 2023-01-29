package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb+srv://user:aabbccdd1234@cluster0.0esyhx6.mongodb.net/?retryWrites=true&w=majority"

var mongoDB *mongo.Database

func GetMongoDB() (*mongo.Database, error) {
	if mongoDB == nil {
		err := initDatabase()
		if err != nil {
			return nil, err
		}
	}

	return mongoDB, nil
}

func initDatabase() error {
	// Create a new client and connect to the server
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	mongoDB = mongoClient.Database("Magna")

	return nil
}
