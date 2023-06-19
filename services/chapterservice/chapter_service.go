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

func GroupMangaToChapter(chapterId []primitive.ObjectID) (map[string][]model.Chapter, error) {
	listItem, err := new(model.Chapter).GetItemListFromObjectIdGroupByManga(chapterId)
	if err != nil {
		return nil, err
	}
	resMap := make(map[string][]model.Chapter)
	for _, item := range listItem {
		resMap[item.Manga.Hex()] = append(resMap[item.Manga.Hex()], item)
	}

	return resMap, nil
}

func CreateChapter(mangaId string, title string, cover string, price uint, images []string) (string, error) {
	mangaObjId, err := primitive.ObjectIDFromHex(mangaId)
	if err != nil {
		return "", err
	}

	chapter := new(model.Chapter)
	chapter.Manga = mangaObjId
	chapter.Name = title
	chapter.Cover = cover
	chapter.Price = price
	chapter.Images = images
	id, err := chapter.InsertToDatabase()
	if err != nil {
		return "", err
	}

	manga := new(model.Manga)
	manga.GetItemFromObjectId(mangaObjId)
	err = manga.UpdateChapter(chapter)
	if err != nil {
		return "", err
	}

	return id.Hex(), nil
}
