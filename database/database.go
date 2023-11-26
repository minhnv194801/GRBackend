package database

import (
	"context"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDB *mongo.Database
var runOnce sync.Once

func GetMongoDB() (*mongo.Database, error) {
	runOnce.Do(func() {
		uri := os.Getenv("DATABASE_URI")
		databaseName := os.Getenv("DATABASE_NAME")

		mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err != nil {
			return
		}

		mongoDB = mongoClient.Database(databaseName)
	})

	return mongoDB, nil
}
