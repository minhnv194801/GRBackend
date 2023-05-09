package homecontroller

import (
	"fmt"
	"log"
	"magna/model"
	"magna/responses"
	"magna/services/mangaservice"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetListHotItems(c *gin.Context) {
	countParam := c.Param("count")
	count, err := strconv.Atoi(countParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad params"})
		return
	}

	hotItemList, err := mangaservice.GetListHotItems(count)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Error in system"})
		return
	}

	fmt.Println(count)
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
	countParam := c.Param("count")
	count, err := strconv.Atoi(countParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad params"})
		return
	}

	recommendList, err := mangaservice.GetListRecommendation(count)
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
	positionParam := c.Param("position")
	countParam := c.Param("count")
	count, err := strconv.Atoi(countParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad params"})
		return
	}
	position, err := strconv.Atoi(positionParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad params"})
		return
	}

	newestList, totalCount, err := mangaservice.GetNewestList(position, count)
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
