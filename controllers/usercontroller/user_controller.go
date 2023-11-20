package usercontroller

import (
	"log"
	"magna/model"
	"magna/requests"
	"magna/responses"
	"magna/services/chapterservice"
	"magna/services/mangaservice"
	"magna/services/reportservice"
	"magna/services/sessionservice"
	"magna/services/userservice"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}
	user, err := userservice.GetUserInfo(userId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
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
		return
	}
	user, err := userservice.GetUserInfo(userId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}

	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.DisplayName = request.Username
	user.Avatar = request.Avatar
	user.Gender = request.Gender
	log.Println("request", request.Gender)
	if user.UpdateInfo() != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Success"})
}

func GetOwnedChapter(c *gin.Context) {
	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}
	user, err := userservice.GetUserInfo(userId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}
	var response []responses.OwnedChapterResponse
	if len(user.OwnedChapters) == 0 {
		c.IndentedJSON(http.StatusOK, response)
		return
	}
	ownedMangaMap, err := chapterservice.GroupMangaToChapter(user.OwnedChapters)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}
	for mangaId, chapterList := range ownedMangaMap {
		var res responses.OwnedChapterResponse
		manga, err := mangaservice.GetMangaInfo(mangaId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
			return
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

func GetFavoriteMangaList(c *gin.Context) {
	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}
	user, err := userservice.GetUserInfo(userId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}

	var responseList []responses.FavoriteItem
	if len(user.FollowMangas) == 0 {
		c.IndentedJSON(http.StatusOK, responseList)
		return
	}
	favoriteList, err := new(model.Manga).GetNewestItemListFromObjectId(user.FollowMangas)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Error in system"})
		return
	}

	for _, item := range favoriteList {
		var response responses.FavoriteItem
		response.Id = item.Id.Hex()
		response.Title = item.Name
		response.Cover = item.Cover
		chapterList, _ := new(model.Chapter).GetMangaNewestChapterList(item.Id, 3)
		for _, chapter := range chapterList {
			var chapterItem responses.FavoriteChapter
			chapterItem.Id = chapter.Id.Hex()
			chapterItem.Title = chapter.Name
			chapterItem.UpdateTime = chapter.UpdateTime
			response.ChapterList = append(response.ChapterList, chapterItem)
		}
		responseList = append(responseList, response)
	}

	c.IndentedJSON(http.StatusOK, responseList)
}

func GetUserReport(c *gin.Context) {
	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}
	reportList, err := reportservice.GetUserReport(userId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}
	var responseList []responses.ReportResponse
	for _, report := range reportList {
		var response responses.ReportResponse
		var chapter model.Chapter
		chapter.GetItemFromObjectId(report.Chapter)
		response.ChapterTitle = chapter.Name
		response.ChapterCover = chapter.Cover
		response.ChapterId = chapter.Id.Hex()
		response.Content = report.Content
		response.Response = report.Response
		response.Status = report.Status
		response.TimeCreated = int(report.TimeCreated)
		responseList = append(responseList, response)
	}

	c.IndentedJSON(http.StatusOK, responseList)
}

func GetUserRecommendation(c *gin.Context) {
	countParam := c.Param("count")
	count, err := strconv.Atoi(countParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad params"})
		return
	}

	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		userId = ""
	}

	mangaList, err := userservice.GetUserRecommendations(userId, count)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}

	var responseList []responses.RecommendedItem
	for _, manga := range mangaList {
		var response responses.RecommendedItem
		response.Id = manga.Id.Hex()
		response.Title = manga.Name
		response.Cover = manga.Cover
		response.Description = manga.Description
		response.Status = int(manga.Status)
		response.Tags = manga.Tags
		var sum float32
		for _, value := range manga.Rated {
			sum += float32(value)
		}
		if len(manga.Rated) != 0 {
			response.Rating = sum / float32(len(manga.Rated))
		}
		responseList = append(responseList, response)
	}

	c.IndentedJSON(http.StatusOK, responseList)
}

func GetTotalCount(c *gin.Context) {
	totalCount, err := userservice.GetTotalCount()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Error in system"})
		return
	}

	c.IndentedJSON(http.StatusOK, totalCount)
}
