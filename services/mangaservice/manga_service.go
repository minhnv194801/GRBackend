package mangaservice

import (
	"errors"
	"fmt"
	"magna/model"
	"sort"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	hotMangaMap map[primitive.ObjectID]int = make(map[primitive.ObjectID]int)
)

func GetItemFromId(mangaId string) (*model.Manga, error) {
	objId, err := primitive.ObjectIDFromHex(mangaId)
	if err != nil {
		return nil, err
	}
	manga := new(model.Manga)
	err = manga.GetItemFromObjectId(objId)
	if err != nil {
		return nil, err
	}
	return manga, nil
}

func GetMangaInfo(mangaId string) (*model.Manga, error) {
	objId, err := primitive.ObjectIDFromHex(mangaId)
	if err != nil {
		return nil, err
	}
	manga := new(model.Manga)
	err = manga.GetItemFromObjectId(objId)
	if err != nil {
		return nil, err
	}
	hotMangaMap[objId]++

	return manga, err
}

func GetListHotItems(count int) ([]model.Manga, error) {
	fmt.Println("Start to serve list hot items")
	fmt.Println(hotMangaMap)
	result := make([]model.Manga, 0)

	objIDs := make([]primitive.ObjectID, 0, len(hotMangaMap))
	for objID := range hotMangaMap {
		objIDs = append(objIDs, objID)
	}
	sort.SliceStable(objIDs, func(i, j int) bool {
		return hotMangaMap[objIDs[i]] > hotMangaMap[objIDs[j]]
	})
	if len(objIDs) < count {
		hotMangaList, err := new(model.Manga).GetItemListFromObjectId(objIDs)
		if err != nil {
			return nil, err
		}
		result = append(result, hotMangaList...)
		remainRandomCount := count - len(objIDs)
		randomList, err := new(model.Manga).GetRandomExcludedItemListFromObjectId(objIDs, remainRandomCount)
		if err != nil {
			return nil, err
		}
		result = append(result, randomList...)
	} else {
		hotMangaList, err := new(model.Manga).GetItemListFromObjectId(objIDs[:count])
		if err != nil {
			return nil, err
		}
		result = append(result, hotMangaList...)
	}

	return result, nil
}

func GetMangaChapterList(mangaId string, position, count int) ([]model.Chapter, error) {
	objId, err := primitive.ObjectIDFromHex(mangaId)
	if err != nil {
		return nil, err
	}

	manga := new(model.Manga)
	err = manga.GetItemFromObjectId(objId)
	if err != nil {
		return nil, err
	}

	if position >= len(manga.Chapters) {
		return nil, errors.New("Invalid position")
	}

	if position+count >= len(manga.Chapters) {
		chapterList, err := new(model.Chapter).GetItemListFromObjectId(manga.Chapters[position:])
		if err != nil {
			return nil, err
		}
		return chapterList, nil
	} else {
		chapterList, err := new(model.Chapter).GetItemListFromObjectId(manga.Chapters[position : position+count])
		if err != nil {
			return nil, err
		}
		return chapterList, nil
	}
}

func GetListRecommendation(count int) ([]model.Manga, error) {
	return new(model.Manga).GetListRecommendManga(count)
}

func GetNewestList(position, count int) ([]model.Manga, error) {
	return new(model.Manga).GetNewestItemList(position, count)
}

func GetTotalCount() (int, error) {
	return new(model.Manga).GetTotalCount()
}

func ClearHotMangaMap() {
	hotMangaMap = make(map[primitive.ObjectID]int)
}
