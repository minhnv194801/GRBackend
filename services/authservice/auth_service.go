package authservice

import (
	"errors"
	"log"
	"magna/model"
	"magna/services/sessionservice"
	"magna/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Login(email, password string) (sskey string, refreshkey string, id string, username, avatar string, err error) {
	if !utils.ValidateEmail(email) {
		return "", "", "", "", "", errors.New("Email không hợp lệ")
	}
	if !utils.ValidatePassword(password) {
		return "", "", "", "", "", errors.New("Password không hợp lệ")
	}
	user := new(model.User)
	err = user.GetItemFromEmail(email)
	if err != nil {
		return "", "", "", "", "", errors.New("Không tồn tại tài khoản với email này")
	}
	err = utils.CheckPasswordHash(user.Password, password)
	if err != nil {
		return "", "", "", "", "", errors.New("Sai mật khẩu")
	}

	id = user.Id.Hex()
	username = user.DisplayName
	avatar = user.Avatar
	sskey, refreshkey, err = sessionservice.CreateSession(id)
	if err != nil {
		return "", "", "", "", "", errors.New("Internal server error")
	}

	return sskey, refreshkey, id, username, avatar, nil
}

func Register(email, password, rePassword string) (sskey string, refreshkey string, id string, username, avatar string, err error) {
	if !utils.ValidateEmail(email) {
		return "", "", "", "", "", errors.New("Email không hợp lệ")
	}
	if !utils.ValidatePassword(password) {
		return "", "", "", "", "", errors.New("Password không hợp lệ")
	}
	if password != rePassword {
		return "", "", "", "", "", errors.New("Mật khẩu nhập lại khác với mật khẩu")
	}
	user := new(model.User)
	_, err = user.CreateNewUser(email, password, "Người dùng")
	if err != nil {
		log.Println(err.Error(), "err.Error() services/userservice/user_service.go:38")
		return "", "", "", "", "", err
	}

	id = user.Id.Hex()
	username = user.DisplayName
	avatar = user.Avatar
	sskey, refreshkey, err = sessionservice.CreateSession(id)
	if err != nil {
		return "", "", "", "", "", errors.New("Internal server error")
	}

	return sskey, refreshkey, id, username, avatar, nil
}

func CheckAdmin(userId string) (bool, error) {
	if userId == "" {
		return false, nil
	}

	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return false, err
	}
	user := new(model.User)
	user.Id = objId
	return user.IsAdmin()
}
