package sessionservice

import (
	"fmt"
	"magna/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtSecretKey = []byte("kcxhrtodjhbf;bc;jtfpfd")
var refreshSecretKey = []byte("vgfkvcjbprsrpsgdfbnjo")

func CreateSession(id string) (sessionkey string, refreshkey string, err error) {
	claims := jwt.MapClaims{}
	claims["id"] = id
	//TODO: put expired time in config
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() //Token hết hạn sau 24 giờ
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	sessionkey, err = token.SignedString(jwtSecretKey)
	if err != nil {
		return "", "", err
	}

	claims["exp"] = time.Now().Add(time.Hour * 48).Unix() //Token hết hạn sau 48 giờ
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshkey, err = token.SignedString(refreshSecretKey)
	if err != nil {
		return "", "", err
	}

	return sessionkey, refreshkey, nil
}

func RefreshSession(key string) (sessionkey string, refreshkey string, userId string, username string, avatar string, err error) {
	token, err := jwt.Parse(key, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		return refreshSecretKey, nil
	})
	if err != nil {
		return "", "", "", "", "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	userId = fmt.Sprintf("%v", claims["id"])
	sessionkey, refreshkey, err = CreateSession(userId)
	if err != nil {
		return "", "", "", "", "", err
	}
	user := new(model.User)
	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return "", "", "", "", "", err
	}
	user.GetItemFromObjectId(objId)
	username = user.DisplayName
	avatar = user.Avatar

	return sessionkey, refreshkey, userId, username, avatar, nil
}

func CheckSession(sessionkey string) (*jwt.Token, error) {
	return jwt.Parse(sessionkey, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		return jwtSecretKey, nil
	})
}
