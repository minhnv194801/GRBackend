package ratingservice

import (
	"magna/services/mangaservice"
	"magna/services/userservice"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetRating(userId, mangaId string, rate int) error {
	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	mangaObjId, err := primitive.ObjectIDFromHex(mangaId)
	if err != nil {
		return err
	}

	manga, err := mangaservice.GetMangaInfo(mangaId)
	if err != nil {
		return err
	}
	user, err := userservice.GetUserInfo(userId)
	if err != nil {
		return err
	}

	if manga.Rated == nil {
		manga.Rated = make(map[primitive.ObjectID]int)
	}
	manga.Rated[mangaObjId] = rate
	if user.Rate == nil {
		user.Rate = make(map[primitive.ObjectID]int)
	}
	user.Rate[userObjId] = rate

	err = manga.SetRated()
	if err != nil {
		return err
	}
	return user.SetRate()
}
