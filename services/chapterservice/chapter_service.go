package chapterservice

import (
	"magna/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetChapterInfo(chapterId string) (*model.Chapter, error) {
	objId, err := primitive.ObjectIDFromHex(chapterId)
	if err != nil {
		return nil, err
	}
	chapter := new(model.Chapter)
	err = chapter.GetItemFromObjectId(objId)
	if err != nil {
		return nil, err
	}
	return chapter, err
}

func GetSameMangaChapterList(chapterId string) ([]model.Chapter, error) {
	objId, err := primitive.ObjectIDFromHex(chapterId)
	if err != nil {
		return nil, err
	}
	var chapter model.Chapter
	err = chapter.GetItemFromObjectId(objId)
	if err != nil {
		return nil, err
	}

	return chapter.GetMangaChapterList(chapter.Manga)
}

func CheckIsOwner(chapterId, userId string) (bool, error) {
	chapterObjId, err := primitive.ObjectIDFromHex(chapterId)
	if err != nil {
		return false, err
	}
	userObjId, _ := primitive.ObjectIDFromHex(userId)
	var chapter model.Chapter
	chapter.Id = chapterObjId
	return chapter.IsOwned(userObjId)
}
