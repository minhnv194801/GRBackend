package tests

import (
	"magna/model"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetListItems(t *testing.T) {
	objIdLst := make([]primitive.ObjectID, 0)
	id, _ := primitive.ObjectIDFromHex("63d76eec1adb5dbe924795a7")
	objIdLst = append(objIdLst, id)
	id, _ = primitive.ObjectIDFromHex("63d7619d1adb5dbe924794cb")
	objIdLst = append(objIdLst, id)
	t.Log(objIdLst)

	itemList, err := new(model.Manga).GetItemListFromObjectId(objIdLst)
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log(len(itemList))
	for _, item := range itemList {
		t.Log(item.Name)
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
	recommended, err := new(model.Manga).GetNewestItemList(8, 9)
	if err != nil {
		t.Error(err.Error())
		return
	}

	for _, recomm := range recommended {
		t.Log(recomm.Name)
	}
}
