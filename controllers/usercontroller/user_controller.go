package usercontroller

import (
	"log"
	"magna/requests"
	"magna/responses"
	"magna/services/chapterservice"
	"magna/services/mangaservice"
	"magna/services/sessionservice"
	"magna/services/userservice"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}
	user, err := userservice.GetUserInfo(userId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	var response responses.UserInfoResponse
	response.Email = user.Email
	response.FirstName = user.FirstName
	response.LastName = user.LastName
	response.Gender = user.Gender
	response.Role = user.Role
	c.IndentedJSON(http.StatusOK, response)
}

func UpdateUserInfo(c *gin.Context) {
	var request requests.UpdateUserInfoRequest

	err := c.BindJSON(&request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}
	user, err := userservice.GetUserInfo(userId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.DisplayName = request.Username
	user.Avatar = request.Avatar
	user.Gender = request.Gender
	log.Println("request", request.Gender)
	if user.UpdateInfo() != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Success"})
}

func GetOwnedChapter(c *gin.Context) {
	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}
	user, err := userservice.GetUserInfo(userId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}
	ownedMangaMap, err := chapterservice.GroupMangaToChapter(user.OwnedChapters)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}
	var response []responses.OwnedChapterResponse
	for mangaId, chapterList := range ownedMangaMap {
		var res responses.OwnedChapterResponse
		manga, err := mangaservice.GetMangaInfo(mangaId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		}
		res.Id = manga.Id.Hex()
		res.Title = manga.Name
		res.Cover = manga.Cover
		for _, chapter := range chapterList {
			var chapterResponse responses.OwnedChapterItem
			chapterResponse.Id = chapter.Id.Hex()
			chapterResponse.Title = chapter.Name
			res.Chapters = append(res.Chapters, chapterResponse)
		}
		response = append(response, res)
	}

	c.IndentedJSON(http.StatusOK, response)
}
