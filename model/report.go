package model

import (
	"context"
	"magna/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (report *Report) GetUserReport(userId primitive.ObjectID) ([]Report, error) {
	coll, err := database.GetReportCollection()
	if err != nil {
		return nil, err
	}

	listItem := make([]Report, 0)
	filter := bson.M{"user": userId}
	opts := options.Find().SetSort(bson.D{{"timeCreated", -1}})
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	err = cursor.All(context.TODO(), &listItem)
	if err != nil {
		return nil, err
	}

	return listItem, nil
}
