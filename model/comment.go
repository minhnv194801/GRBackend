package model

import (
	"context"
	"magna/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Comment struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Manga       primitive.ObjectID `bson:"manga"`
	User        primitive.ObjectID `bson:"user"`
	Content     string             `bson:"content"`
	TimeCreated uint               `bson:"timeCreated"`
}

func (comment *Comment) InsertToDatabase() (primitive.ObjectID, error) {
	coll, err := database.GetCommentCollection()
	if err != nil {
		return [12]byte{}, err
	}

	result, err := coll.InsertOne(context.TODO(), comment)
	if err != nil {
		return [12]byte{}, err
	}

	comment.Id = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID.(primitive.ObjectID), nil
}

func (comment *Comment) GetListCommentFromMangaId(mangaId primitive.ObjectID, position, count int) ([]Comment, error) {
	coll, err := database.GetCommentCollection()
	if err != nil {
		return nil, err
	}

	filter := bson.D{primitive.E{Key: "manga", Value: mangaId}}
	opts := options.Find().SetSort(bson.D{{"timeCreated", -1}}).SetSkip(int64(position)).SetLimit(int64(count))
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	listItem := make([]Comment, 0)
	err = cursor.All(context.TODO(), &listItem)
	if err != nil {
		return nil, err
	}

	return listItem, err
}

func (comment *Comment) GetMangaCommentCount(mangaId primitive.ObjectID) (int, error) {
	coll, err := database.GetCommentCollection()
	if err != nil {
		return 0, err
	}

	filter := bson.D{primitive.E{Key: "manga", Value: mangaId}}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (comment *Comment) CreateNewComment() (primitive.ObjectID, error) {
	comment.TimeCreated = uint(time.Now().Unix())

	coll, err := database.GetCommentCollection()
	if err != nil {
		return [12]byte{}, err
	}

	result, err := coll.InsertOne(context.TODO(), comment)
	if err != nil {
		return [12]byte{}, err
	}

	comment.Id = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID.(primitive.ObjectID), nil
}
