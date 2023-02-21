package homecontroller

import (
	"magna/model"
	"magna/requests"
	"magna/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetListRecommendation(c *gin.Context) {
	var request requests.RecommendListRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	recommendList, err := new(model.Manga).GetListRecommendManga(request.Count)
	if err != nil {
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

	newestList, err := new(model.Manga).GetNewestItemList(request.Postition, request.Count)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Error in system"})
		return
	}

	var responseList []responses.RecommendResponse
	for _, item := range newestList {
		var response responses.RecommendResponse
		response.Id = item.Id.Hex()
		response.Title = item.Name
		response.Image = item.Cover
		responseList = append(responseList, response)
	}

	c.IndentedJSON(http.StatusOK, responseList)
}
