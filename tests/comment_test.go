package tests

import (
	"magna/model"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetListCommentFromMangaId(t *testing.T) {
	objId, err := primitive.ObjectIDFromHex("63d7619d1adb5dbe924794cb")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("OK")
	t.Log(new(model.Comment).GetListCommentFromMangaId(objId, 0, 10))
}
