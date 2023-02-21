package homecontroller

import (
	"magna/model"
	"magna/requests"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetChapterInfo(c *gin.Context) {
	var request requests.RecommendListRequest
	c.BindJSON(&request)

	recommendList, err := (*model.Manga).GetListRecommendManga(request.Count)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Error in system"})
	}

	// var response responses.ReadResponse

	c.IndentedJSON(http.StatusOK, recommendList)
}
