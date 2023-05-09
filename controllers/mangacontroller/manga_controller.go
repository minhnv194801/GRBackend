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
	"magna/services/ratingservice"
	"magna/services/sessionservice"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetMangaInfo(c *gin.Context) {
	id := c.Param("mangaid")
	log.Println("MangaId:", id)
	sessionkey := c.GetHeader("Authorization")
	userId, _ := sessionservice.ExtractUserIdFromSessionKey(sessionkey)

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
	isFavorite, _ := mangaservice.CheckIsFavorite(id, userId)
	response.IsFavorite = isFavorite
	objId, err := primitive.ObjectIDFromHex(userId)
	if manga.Rated == nil {
		manga.Rated = make(map[primitive.ObjectID]int)
	}
	fmt.Println(manga.Rated[objId])
	if err == nil {
		response.UserRating = uint(manga.Rated[objId])
	} else {
		response.UserRating = 0
	}
	var sum float32
	for _, value := range manga.Rated {
		sum += float32(value)
	}
	if len(manga.Rated) != 0 {
		response.AvgRating = sum / float32(len(manga.Rated))
	}
	response.RatingCount = uint(len(manga.Rated))
	response.Description = manga.Description
	response.ChapterCount = len(manga.Chapters)

	c.IndentedJSON(http.StatusOK, response)

}

func GetMangaChapterList(c *gin.Context) {
	id := c.Param("mangaid")
	sessionkey := c.GetHeader("Authorization")
	userId, _ := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
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

	var responseList []responses.ChapterListResponse
	chapterList, err := mangaservice.GetMangaChapterList(id, position, count)
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
		isOwned, _ := chapterservice.CheckIsOwner(chapter.Id.Hex(), userId)
		response.IsOwned = isOwned
		response.UpdateTime = chapter.UpdateTime
		responseList = append(responseList, response)
	}

	c.IndentedJSON(http.StatusOK, responseList)
}

func GetCommentList(c *gin.Context) {
	id := c.Param("mangaid")
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

	var responseList responses.CommentListResponse
	commentList, err := commentservice.GetCommentListFromMangaId(id, position, count)
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

func SetRating(c *gin.Context) {
	id := c.Param("mangaid")
	var request requests.RatingRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}
	sessionkey := c.GetHeader("Authorization")
	userId, err := sessionservice.ExtractUserIdFromSessionKey(sessionkey)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	err = ratingservice.SetRating(userId, id, request.Rating)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Success"})
}
