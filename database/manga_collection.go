package database

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

var mangaCollection *mongo.Collection
var mangaRunOnce sync.Once

func GetMangaCollection() (*mongo.Collection, error) {
	db, err := GetMongoDB()
	if err != nil {
		return nil, err
	}
	mangaCollection = db.Collection("Manga")

	return mangaCollection, nil
}
