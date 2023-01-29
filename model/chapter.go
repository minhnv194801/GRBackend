package model

import (
	"context"
	"magna/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Chapter struct {
	Id         primitive.ObjectID   `bson:"_id,omitempty"`
	Manga      primitive.ObjectID   `bson:"manga"`
	Name       string               `bson:"name"`
	Cover      string               `bson:"cover"`
	Price      uint                 `bson:"price"`
	UpdateTime uint                 `bson:"updateTime"`
	Images     []string             `bson:"images"`
	OwnedUsers []primitive.ObjectID `bson:"ownedUsers"`
	Reports    []primitive.ObjectID `bson:"reports"`
}

func (chapter *Chapter) InsertToDatabase() (primitive.ObjectID, error) {
	db, err := database.GetMongoDB()
	if err != nil {
		return [12]byte{}, err
	}
	coll := db.Collection("Chapter")

	id, err := getExistedChapterID(chapter.Manga, chapter.Name)
	if id != primitive.NilObjectID {
		chapter.Id = id
		return id, err
	}

	result, err := coll.InsertOne(context.TODO(), chapter)
	if err != nil {
		return [12]byte{}, err
	}

	chapter.Id = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID.(primitive.ObjectID), nil
}

func getExistedChapterID(mangaId primitive.ObjectID, name string) (primitive.ObjectID, error) {
	db, err := database.GetMongoDB()
	if err != nil {
		return [12]byte{}, err
	}
	coll := db.Collection("Chapter")

	var doc Chapter
	filter := bson.D{{"manga", mangaId}, {"name", name}}
	opts := options.FindOne()
	found := coll.FindOne(context.TODO(), filter, opts)
	if found.Err() != nil {
		return [12]byte{}, found.Err()
	}
	err = found.Decode(&doc)
	if err != nil {
		return [12]byte{}, err
	}
	return doc.Id, nil
}
