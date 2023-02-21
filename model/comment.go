package model

import (
	"context"
	"magna/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
