package database

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

var chapterCollection *mongo.Collection
var chapterRunOnce sync.Once

func GetChapterCollection() (*mongo.Collection, error) {
	chapterRunOnce.Do(func() {
		db, err := GetMongoDB()
		if err != nil {
			return
		}
		chapterCollection = db.Collection("Chapter")
	})

	return chapterCollection, nil
}
