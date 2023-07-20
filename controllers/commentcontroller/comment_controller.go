package commentcontroller

import (
	"magna/requests"
	"magna/services/commentservice"
	"magna/services/sessionservice"
	"magna/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateNewComment(c *gin.Context) {
	var request requests.CreateCommentRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	if utils.CheckEmptyString(request.Content) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Empty comment content"})
		return
	}
	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	mangaId := c.Param("mangaid")
	err = commentservice.CreateNewComment(userId, mangaId, request.Content)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Success"})
}

func GetTotalCount(c *gin.Context) {
	totalCount, err := commentservice.GetTotalCount()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Error in system"})
		return
	}

	c.IndentedJSON(http.StatusOK, totalCount)
}
