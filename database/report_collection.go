package database

import "go.mongodb.org/mongo-driver/mongo"

func GetReportCollection() (*mongo.Collection, error) {
	db, err := GetMongoDB()
	if err != nil {
		return nil, err
	}
	collection := db.Collection("Reports")

	return collection, nil
}
