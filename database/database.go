package database

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDB *mongo.Database
var runOnce sync.Once

const uri = "mongodb+srv://user:aabbccdd1234@cluster0.0esyhx6.mongodb.net/?retryWrites=true&w=majority"

func GetMongoDB() (*mongo.Database, error) {
	runOnce.Do(func() {
		mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err != nil {
			return
		}

		mongoDB = mongoClient.Database("Magna-test")
	})

	return mongoDB, nil
}
