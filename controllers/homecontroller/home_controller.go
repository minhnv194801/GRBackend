package homecontroller

import (
	"log"
	"magna/model"
	"magna/requests"
	"magna/responses"
	"magna/services/mangaservice"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetListHotItems(c *gin.Context) {
	var request requests.HotItemsListRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	hotItemList, err := mangaservice.GetListHotItems(request.Count)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Error in system"})
		return
	}

	var responseList []responses.HotItemsResponse
	for _, hotItem := range hotItemList {
		var response responses.HotItemsResponse
		response.Id = hotItem.Id.Hex()
		response.Title = hotItem.Name
		response.Image = hotItem.Cover
		responseList = append(responseList, response)
	}

	c.IndentedJSON(http.StatusOK, responseList)
}

func GetListRecommendation(c *gin.Context) {
	var request requests.RecommendListRequest
	err := c.BindJSON(&request)
	if err != nil {
		log.Println("ERROR:", err.Error(), "home_controller GetListRecommendation")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	recommendList, err := mangaservice.GetListRecommendation(request.Count)
	if err != nil {
		log.Println("ERROR:", err.Error(), "home_controller GetListRecommendation")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Error in system"})
		return
	}

	var responseList []responses.RecommendResponse
	for _, recommend := range recommendList {
		var response responses.RecommendResponse
		response.Id = recommend.Id.Hex()
		response.Title = recommend.Name
		response.Image = recommend.Cover
		responseList = append(responseList, response)
	}

	c.IndentedJSON(http.StatusOK, responseList)
}

func GetNewestList(c *gin.Context) {
	var request requests.NewestListRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	newestList, totalCount, err := mangaservice.GetNewestList(request.Postition, request.Count)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Error in system"})
		return
	}

	var responseList responses.NewestResponse
	for _, item := range newestList {
		var response responses.NewestItem
		response.Id = item.Id.Hex()
		response.Title = item.Name
		response.Cover = item.Cover
		chapterList, _ := new(model.Chapter).GetMangaNewestChapterList(item.Id, 3)
		for _, chapter := range chapterList {
			var chapterItem responses.NewestChapter
			chapterItem.Id = chapter.Id.Hex()
			chapterItem.Name = chapter.Name
			chapterItem.UpdateTime = chapter.UpdateTime
			response.ChapterList = append(response.ChapterList, chapterItem)
		}
		responseList.Data = append(responseList.Data, response)
	}
	responseList.TotalCount = totalCount

	c.IndentedJSON(http.StatusOK, responseList)
}

func GetTotalCount(c *gin.Context) {
	totalCount, err := mangaservice.GetTotalCount()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Error in system"})
		return
	}

	c.IndentedJSON(http.StatusOK, totalCount)
}
