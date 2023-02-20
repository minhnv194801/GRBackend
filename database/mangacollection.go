package database

import "go.mongodb.org/mongo-driver/mongo"

func GetMangaCollection() (*mongo.Collection, error) {
	db, err := GetMongoDB()
	if err != nil {
		return nil, err
	}
	collection := db.Collection("Manga")

	return collection, nil
}
