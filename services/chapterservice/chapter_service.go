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
