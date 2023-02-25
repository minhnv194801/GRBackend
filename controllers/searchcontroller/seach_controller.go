package searchcontroller

import (
	"magna/model"
	"magna/requests"
	"magna/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	var request requests.SearchRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	var response responses.SearchResponse
	itemList, totalCount, err := new(model.Manga).Filter(request.Query, request.Tags, request.Position, request.Count)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}

	for _, item := range itemList {
		var result responses.SearchItem
		result.Id = item.Id.Hex()
		result.Cover = item.Cover
		result.Title = item.Name
		chapterList, _ := new(model.Chapter).GetMangaNewestChapterList(item.Id, 3)
		for _, chapter := range chapterList {
			var chapterItem responses.NewestChapter
			chapterItem.Id = chapter.Id.Hex()
			chapterItem.Name = chapter.Name
			chapterItem.UpdateTime = chapter.UpdateTime
			result.ChapterList = append(result.ChapterList, chapterItem)
		}
		response.Data = append(response.Data, result)
	}
	response.TotalCount = totalCount

	c.IndentedJSON(http.StatusOK, response)
}
