package model

import (
	"context"
	"magna/database"
	"magna/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Report struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Chapter     primitive.ObjectID `bson:"chapter" json:"chapter"`
	User        primitive.ObjectID `bson:"user" json:"user"`
	Content     string             `bson:"content" json:"content"`
	TimeCreated uint               `bson:"timeCreated" json:"timeCreated"`
	Status      int                `bson:"status" json:"status"`
	Response    string             `bson:"response" json:"response"`
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

	user := new(User)
	err = user.GetItemFromObjectId(report.User)
	if err != nil {
		return [12]byte{}, err
	}
	err = user.AddReport(report.Id)
	if err != nil {
		return [12]byte{}, err
	}

	chapter := new(Chapter)
	err = chapter.GetItemFromObjectId(report.Chapter)
	if err != nil {
		return [12]byte{}, err
	}
	err = chapter.AddReport(report.Id)
	if err != nil {
		return [12]byte{}, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (report *Report) GetItemList(position, count int, sortField, sortType string) ([]Report, int, error) {
	coll, err := database.GetReportCollection()
	if err != nil {
		return nil, 0, err
	}

	listItem := make([]Report, 0)
	filter := bson.D{{}}
	opts := options.Find()
	opts.SetSkip(int64(position))
	if sortField == "id" {
		sortField = "_id"
	}
	if sortType == "ASC" {
		opts.SetSort(bson.M{utils.FirstLetterToLower(sortField): 1})
	} else {
		opts.SetSort(bson.M{utils.FirstLetterToLower(sortField): -1})
	}
	opts.SetLimit(int64(count))

	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())
	err = cursor.All(context.TODO(), &listItem)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}

	if count < len(listItem) {
		return listItem[:count], int(totalCount), nil
	} else {
		return listItem[:], int(totalCount), nil
	}
}

func (report *Report) GetItemListFilterByUser(position, count int, sortField, sortType, filterValue string) ([]Report, int, error) {
	coll, err := database.GetReportCollection()
	if err != nil {
		return nil, 0, err
	}
	filterValueObjId, err := primitive.ObjectIDFromHex(filterValue)
	if err != nil {
		return nil, 0, err
	}

	listItem := make([]Report, 0)
	filter := bson.M{"user": filterValueObjId}
	opts := options.Find()
	opts.SetSkip(int64(position))
	if sortField == "id" {
		sortField = "_id"
	}
	if sortType == "ASC" {
		opts.SetSort(bson.M{utils.FirstLetterToLower(sortField): 1})
	} else {
		opts.SetSort(bson.M{utils.FirstLetterToLower(sortField): -1})
	}
	opts.SetLimit(int64(count))

	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())
	err = cursor.All(context.TODO(), &listItem)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}

	if count < len(listItem) {
		return listItem[:count], int(totalCount), nil
	} else {
		return listItem[:], int(totalCount), nil
	}
}

func (report *Report) GetItemListFilterByChapter(position, count int, sortField, sortType, filterValue string) ([]Report, int, error) {
	coll, err := database.GetReportCollection()
	if err != nil {
		return nil, 0, err
	}
	filterValueObjId, err := primitive.ObjectIDFromHex(filterValue)
	if err != nil {
		return nil, 0, err
	}

	listItem := make([]Report, 0)
	filter := bson.M{"chapter": filterValueObjId}
	opts := options.Find()
	opts.SetSkip(int64(position))
	if sortField == "id" {
		sortField = "_id"
	}
	if sortType == "ASC" {
		opts.SetSort(bson.M{utils.FirstLetterToLower(sortField): 1})
	} else {
		opts.SetSort(bson.M{utils.FirstLetterToLower(sortField): -1})
	}
	opts.SetLimit(int64(count))

	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())
	err = cursor.All(context.TODO(), &listItem)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}

	if count < len(listItem) {
		return listItem[:count], int(totalCount), nil
	} else {
		return listItem[:], int(totalCount), nil
	}
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

func (report *Report) GetItemFromObjectId(objID primitive.ObjectID) error {
	coll, err := database.GetReportCollection()
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	opts := options.FindOne()
	found := coll.FindOne(context.TODO(), filter, opts)
	if found.Err() != nil {
		return found.Err()
	}
	err = found.Decode(report)
	if err != nil {
		return err
	}

	return nil
}

func (report *Report) DeleteReportById(id primitive.ObjectID) error {
	err := report.GetItemFromObjectId(id)
	if err != nil {
		return err
	}

	coll, err := database.GetReportCollection()
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", report.Id}}
	_, err = coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	user := new(User)
	err = user.GetItemFromObjectId(report.User)
	// If successful get user then remove report from user
	if err == nil {
		user.RemoveReport(report.Id)
	}

	chapter := new(Chapter)
	err = chapter.GetItemFromObjectId(report.Chapter)
	// If successful get chapter then remove report from chapter
	if err == nil {
		chapter.RemoveReport(report.Id)
	}

	return nil
}

func (report *Report) Respond(respond string) error {
	coll, err := database.GetReportCollection()
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", report.Id}}
	update := bson.D{{"$set", bson.D{
		{"response", respond},
		{"status", 1},
	}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (report *Report) GetTotalCount() (int, error) {
	coll, err := database.GetReportCollection()
	if err != nil {
		return 0, err
	}

	filter := bson.D{{}}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
