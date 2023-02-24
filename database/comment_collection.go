package database

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

var commentCollection *mongo.Collection
var commentRunOnce sync.Once

func GetCommentCollection() (*mongo.Collection, error) {
	commentRunOnce.Do(func() {
		db, err := GetMongoDB()
		if err != nil {
			return
		}
		commentCollection = db.Collection("Comments")
	})

	return commentCollection, nil
}
