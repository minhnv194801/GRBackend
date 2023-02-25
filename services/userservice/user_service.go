package userservice

import (
	"magna/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserInfo(userId string) (*model.User, error) {
	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	user := new(model.User)
	err = user.GetItemFromObjectId(objId)
	if err != nil {
		return nil, err
	}
	return user, err
}
