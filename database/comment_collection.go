package database

import "go.mongodb.org/mongo-driver/mongo"

func GetCommentCollection() (*mongo.Collection, error) {
	db, err := GetMongoDB()
	if err != nil {
		return nil, err
	}
	collection := db.Collection("Comments")

	return collection, nil
}
