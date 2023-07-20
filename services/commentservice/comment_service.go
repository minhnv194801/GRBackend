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

func CreateNewComment(userId, mangaId, commentContent string) error {
	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	mangaObjId, err := primitive.ObjectIDFromHex(mangaId)
	if err != nil {
		return err
	}
	comment := new(model.Comment)
	comment.Manga = mangaObjId
	comment.User = userObjId
	comment.Content = commentContent

	_, err = comment.CreateNewComment()

	return err
}

func GetTotalCount() (int, error) {
	return new(model.Comment).GetTotalCount()
}
