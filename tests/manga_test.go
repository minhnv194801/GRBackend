package tests

import (
	"magna/model"
	"testing"
)

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
