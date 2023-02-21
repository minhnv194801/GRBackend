package model

import (
	"context"
	"magna/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Email         string               `bson:"email"`
	Password      string               `bson:"password"`
	Role          string               `bson:"role"`
	DisplayName   string               `bson:"displayName"`
	Avatar        string               `bson:"avatar"`
	FirstName     string               `bson:"firstName"`
	LastName      string               `bson:"lastName"`
	Gender        int                  `bson:"gender"`
	FollowMangas  []primitive.ObjectID `bson:"followMangas"`
	OwnedChapters []primitive.ObjectID `bson:"ownedChapters"`
	Comments      []primitive.ObjectID `bson:"comments"`
	Reports       []primitive.ObjectID `bson:"reports"`
}

func (user *User) InsertToDatabase() (primitive.ObjectID, error) {
	coll, err := database.GetUserCollection()
	if err != nil {
		return [12]byte{}, err
	}

	result, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		return [12]byte{}, err
	}

	user.Id = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID.(primitive.ObjectID), nil
}
