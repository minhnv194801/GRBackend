package database

import "go.mongodb.org/mongo-driver/mongo"

func GetChapterCollection() (*mongo.Collection, error) {
	db, err := GetMongoDB()
	if err != nil {
		return nil, err
	}
	collection := db.Collection("Chapter")

	return collection, nil
}
