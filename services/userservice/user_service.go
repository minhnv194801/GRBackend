package userservice

import (
	"errors"
	"log"
	"magna/model"
	"magna/utils"

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

func CreateAccount(email, password, role string) (string, error) {
	if !utils.ValidateEmail(email) {
		return "", errors.New("Email không hợp lệ")
	}
	if !utils.ValidatePassword(password) {
		return "", errors.New("Password không hợp lệ")
	}

	user := new(model.User)
	_, err := user.CreateNewUser(email, password, role)
	if err != nil {
		log.Println(err.Error(), "err.Error() services/userservice/user_service.go:38")
		return "", err
	}

	id := user.Id.Hex()

	return id, nil
}

func GetTotalCount() (int, error) {
	return new(model.User).GetTotalCount()
}
