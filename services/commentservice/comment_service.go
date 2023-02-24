package commentservice

import (
	"magna/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCommentListFromMangaId(mangaId string, position, count int) ([]model.Comment, error) {
	objId, err := primitive.ObjectIDFromHex(mangaId)
	if err != nil {
		return nil, err
	}
	return new(model.Comment).GetListCommentFromMangaId(objId, position, count)
}

func GetMangaCommentCount(mangaId string) (int, error) {
	objId, err := primitive.ObjectIDFromHex(mangaId)
	if err != nil {
		return 0, err
	}
	return new(model.Comment).GetMangaCommentCount(objId)
}
