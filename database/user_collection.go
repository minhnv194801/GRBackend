package database

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection
var userRunOnce sync.Once

func GetUserCollection() (*mongo.Collection, error) {
	userRunOnce.Do(func() {
		db, err := GetMongoDB()
		if err != nil {
			return
		}
		userCollection = db.Collection("User")
	})

	return userCollection, nil
}
