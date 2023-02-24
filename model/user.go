package model

import (
	"context"
	"errors"
	"magna/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (user *User) GetItemFromObjectId(objID primitive.ObjectID) error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	opts := options.FindOne()
	found := coll.FindOne(context.TODO(), filter, opts)
	if found.Err() != nil {
		return found.Err()
	}
	err = found.Decode(user)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) GetItemFromEmail(email string) error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "email", Value: email}}
	opts := options.FindOne()
	found := coll.FindOne(context.TODO(), filter, opts)
	if found.Err() != nil {
		return found.Err()
	}
	err = found.Decode(user)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) CreateNewUser() (primitive.ObjectID, error) {
	existed, err := checkExistedEmail(user.Email)
	if err != nil {
		return [12]byte{}, err
	}
	if existed {
		return [12]byte{}, errors.New("Email đã tồn tại")
	}

	user.Role = "Người dùng"
	user.Avatar = "https://st3.depositphotos.com/1767687/16607/v/450/depositphotos_166074422-stock-illustration-default-avatar-profile-icon-grey.jpg"
	user.DisplayName = user.Email
	user.FirstName = "Tên"
	user.LastName = "Họ"
	user.Gender = 0

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

func checkExistedEmail(email string) (bool, error) {
	coll, err := database.GetUserCollection()
	if err != nil {
		return false, err
	}

	filter := bson.D{primitive.E{Key: "email", Value: email}}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	if count != 0 {
		return true, nil
	}
	return false, nil
}
