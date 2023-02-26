package model

import (
	"context"
	"magna/database"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Chapter     primitive.ObjectID `bson:"chapter"`
	User        primitive.ObjectID `bson:"user"`
	Content     string             `bson:"content"`
	TimeCreated uint               `bson:"timeCreated"`
	Status      int                `bson:"status"`
	Response    string             `bson:"response"`
}

func (report *Report) InsertToDatabase() (primitive.ObjectID, error) {
	coll, err := database.GetReportCollection()
	if err != nil {
		return [12]byte{}, err
	}

	result, err := coll.InsertOne(context.TODO(), report)
	if err != nil {
		return [12]byte{}, err
	}

	report.Id = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID.(primitive.ObjectID), nil
}

func (report *Report) CreateNewReport() (primitive.ObjectID, error) {
	report.TimeCreated = uint(time.Now().Unix())
	report.Status = 0
	report.Response = ""

	coll, err := database.GetReportCollection()
	if err != nil {
		return [12]byte{}, err
	}

	result, err := coll.InsertOne(context.TODO(), report)
	if err != nil {
		return [12]byte{}, err
	}

	report.Id = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID.(primitive.ObjectID), nil
}
