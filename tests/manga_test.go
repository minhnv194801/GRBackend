package tests

import (
	"magna/model"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetListItems(t *testing.T) {
	objIdLst := make([]primitive.ObjectID, 0)
	id, _ := primitive.ObjectIDFromHex("63fb22505854cfa4e4ca8cdf")
	objIdLst = append(objIdLst, id)
	// id, _ = primitive.ObjectIDFromHex("63d7619d1adb5dbe924794cb")
	// objIdLst = append(objIdLst, id)
	t.Log(objIdLst)

	itemList, err := new(model.Manga).GetItemListFromObjectId(objIdLst)
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log(len(itemList))
	for _, item := range itemList {
		t.Log(item.Name)
		t.Log(item.Rated)
	}
}

func TestGetRecommendItems(t *testing.T) {
	recommended, err := new(model.Manga).GetListRecommendManga(15)
	if err != nil {
		t.Error(err.Error())
		return
	}

	for _, recomm := range recommended {
		t.Log(recomm.Name)
	}
}

func TestGetNewestItemList(t *testing.T) {
	recommended, totalCount, err := new(model.Manga).GetNewestItemList(8, 9)
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log(totalCount)
	for _, recomm := range recommended {
		t.Log(recomm.Name)
	}
}

func TestGetTotalCount(t *testing.T) {
	count, err := new(model.Manga).GetTotalCount()
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log("Count:", count)
}

func TestFilter(t *testing.T) {
	filtered, totalCount, err := new(model.Manga).Filter("Thanh", []string{}, 0, 9)
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log(totalCount)
	for _, item := range filtered {
		t.Log(item.Name)
	}
}

func TestNewField(t *testing.T) {
	manga := new(model.Manga)
	objId, _ := primitive.ObjectIDFromHex("63f853a88b0c09f7346b721d")
	manga.Rated = make(map[primitive.ObjectID]int)
	manga.Rated[objId] = 5
	manga.UpdateTime = uint(time.Now().Unix())
	manga.InsertToDatabase()
}
