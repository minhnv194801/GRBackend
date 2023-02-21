package database

import "go.mongodb.org/mongo-driver/mongo"

func GetUserCollection() (*mongo.Collection, error) {
	db, err := GetMongoDB()
	if err != nil {
		return nil, err
	}
	collection := db.Collection("User")

	return collection, nil
}
