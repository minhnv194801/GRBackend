package model

import (
	"context"
	"magna/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Manga struct {
	Id            primitive.ObjectID   `bson:"_id,omitempty"`
	Name          string               `bson:"name"`
	AlternateName []string             `bson:"alternateName"`
	Author        []string             `bson:"author"`
	Cover         string               `bson:"cover"`
	Description   string               `bson:"description"`
	Status        Status               `bson:"status"`
	IsRecommended bool                 `bson:"isRecommended"`
	Tags          []string             `bson:"tags"`
	FollowedUsers []primitive.ObjectID `bson:"followedUsers"`
	Chapters      []primitive.ObjectID `bson:"chapters"`
	Comments      []primitive.ObjectID `bson:"comments"`
}

type Status int

const (
	Ongoing Status = iota
	Finished
)

func (manga *Manga) InsertToDatabase() (primitive.ObjectID, error) {
	db, err := database.GetMongoDB()
	if err != nil {
		return [12]byte{}, err
	}
	coll := db.Collection("Manga")

	id, err := getExistedTitleID(manga.Name)
	if id != primitive.NilObjectID {
		manga.Id = id
		return id, err
	}

	result, err := coll.InsertOne(context.TODO(), manga)
	if err != nil {
		return [12]byte{}, err
	}

	manga.Id = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID.(primitive.ObjectID), nil
}

func (manga *Manga) UpdateChapter(chapter *Chapter) error {
	db, err := database.GetMongoDB()
	if err != nil {
		return err
	}
	coll := db.Collection("Manga")

	manga.Chapters = append(manga.Chapters, chapter.Id)
	filter := bson.D{{"_id", manga.Id}}
	update := bson.D{{"$set", bson.D{{"chapters", manga.Chapters}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (manga *Manga) GetMangaFromName(name string) error {
	db, err := database.GetMongoDB()
	if err != nil {
		return err
	}
	coll := db.Collection("Manga")

	var doc Manga
	filter := bson.D{primitive.E{Key: "name", Value: name}}
	opts := options.FindOne()
	found := coll.FindOne(context.TODO(), filter, opts)
	if found.Err() != nil {
		return found.Err()
	}
	err = found.Decode(&doc)
	if err != nil {
		return err
	}
	manga = &doc
	return nil
}

func getExistedTitleID(name string) (primitive.ObjectID, error) {
	db, err := database.GetMongoDB()
	if err != nil {
		return [12]byte{}, err
	}
	coll := db.Collection("Manga")

	var doc Manga
	filter := bson.D{primitive.E{Key: "name", Value: name}}
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
