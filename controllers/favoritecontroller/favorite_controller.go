package favoritecontroller

import (
	"magna/services/mangaservice"
	"magna/services/sessionservice"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetFavorite(c *gin.Context) {
	id := c.Param("mangaid")
	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
	}
	err = mangaservice.SetUserFavorite(userId, id)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Success"})
}
