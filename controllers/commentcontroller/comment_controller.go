package commentcontroller

import (
	"magna/requests"
	"magna/services/commentservice"
	"magna/services/sessionservice"
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
	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
	}
	mangaId := c.Param("mangaId")
	err = commentservice.CreateNewComment(userId, mangaId, request.Content)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Success"})
}
