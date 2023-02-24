package database

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

var reportCollection *mongo.Collection
var reportRunOnce sync.Once

func GetReportCollection() (*mongo.Collection, error) {
	reportRunOnce.Do(func() {
		db, err := GetMongoDB()
		if err != nil {
			return
		}
		reportCollection = db.Collection("Reports")
	})
	return reportCollection, nil
}
