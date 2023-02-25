package mangacontroller

import (
	"fmt"
	"log"
	"magna/model"
	"magna/requests"
	"magna/responses"
	"magna/services/chapterservice"
	"magna/services/commentservice"
	"magna/services/mangaservice"
	"magna/services/sessionservice"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMangaInfo(c *gin.Context) {
	id := c.Param("mangaid")
	log.Println("MangaId:", id)
	//TODO: get user id from token?
	sessionkey := c.GetHeader("Authorization")
	userId, _ := sessionservice.ExtractUserIdFromSessionKey(sessionkey)

	// fmt.Println(c.Request.Header["Authorization"])

	var response responses.MangaInfoResponse

	manga, err := mangaservice.GetMangaInfo(id)
	if err != nil {
		log.Println(err.Error(), "err.Error() GetMangaInfo controllers/readcontrollers/readcontrollers.go:26")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	response.Title = manga.Name
	response.Cover = manga.Cover
	if len(manga.Author) != 0 {
		response.Author = manga.Author[0]
	} else {
		response.Author = ""
	}
	response.Status = uint(manga.Status)
	response.Tags = manga.Tags
	// TODO: check if user already favorite manga
	isFavorite, _ := mangaservice.CheckIsFavorite(id, userId)
	response.IsFavorite = isFavorite
	// TODO: rating services
	response.UserRating = 0
	response.AvgRating = 0
	response.RatingCount = 0
	response.Description = manga.Description
	response.ChapterCount = len(manga.Chapters)

	c.IndentedJSON(http.StatusOK, response)

}

func GetMangaChapterList(c *gin.Context) {
	id := c.Param("mangaid")
	sessionkey := c.GetHeader("Authorization")
	userId, _ := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	var request requests.MangaChapterListRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	//TODO: get user id from token?
	fmt.Println(c.GetHeader("Authorization"))

	var responseList []responses.ChapterListResponse
	chapterList, err := mangaservice.GetMangaChapterList(id, request.Postition, request.Count)
	if err != nil {
		log.Println(err.Error(), "err.Error() GetMangaInfo controllers/readcontrollers/readcontrollers.go:26")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	for _, chapter := range chapterList {
		var response responses.ChapterListResponse
		response.Id = chapter.Id.Hex()
		response.Title = chapter.Name
		response.Cover = chapter.Cover
		response.Price = chapter.Price
		// TODO: Check user ownership
		isOwned, _ := chapterservice.CheckIsOwner(chapter.Id.Hex(), userId)
		response.IsOwned = isOwned
		response.UpdateTime = chapter.UpdateTime
		responseList = append(responseList, response)
	}

	c.IndentedJSON(http.StatusOK, responseList)
}

func GetCommentList(c *gin.Context) {
	id := c.Param("mangaid")
	var request requests.MangaChapterListRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	var responseList responses.CommentListResponse
	commentList, err := commentservice.GetCommentListFromMangaId(id, request.Postition, request.Count)
	if err != nil {
		log.Println(err.Error(), "err.Error() GetCommentList controllers/mangacontrollers/mangacontrollers.go:93")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	for _, comment := range commentList {
		var response responses.CommentInfo
		response.Content = comment.Content
		response.UpdateTime = comment.TimeCreated
		user := new(model.User)
		err := user.GetItemFromObjectId(comment.User)
		if err != nil {
			log.Println(err.Error(), "err.Error() GetCommentList controllers/mangacontrollers/mangacontrollers.go:105")
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		response.Avatar = user.Avatar
		response.Username = user.DisplayName
		responseList.Data = append(responseList.Data, response)
	}
	totalCount, err := commentservice.GetMangaCommentCount(id)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	responseList.TotalCount = totalCount

	c.IndentedJSON(http.StatusOK, responseList)
}

func GetMangaCommentCount(c *gin.Context) {
	id := c.Param("mangaid")

	totalCount, err := commentservice.GetMangaCommentCount(id)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, totalCount)
}
