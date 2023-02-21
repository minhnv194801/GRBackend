package readcontroller

import (
	"fmt"
	"log"
	"magna/model"
	"magna/responses"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetChapterInfo(c *gin.Context) {
	id := c.Param("chapterid")
	log.Println("ChapterId:", id)
	//TODO: check authorization
	fmt.Println(c.GetHeader("Authorization"))
	// fmt.Println(c.Request.Header["Authorization"])

	var response responses.ReadResponse
	var chapter model.Chapter

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err.Error(), "err.Error() GetChapterInfo controllers/readcontrollers/readcontrollers.go:26")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = chapter.GetItemFromObjectId(objId)
	if err != nil {
		log.Println(err.Error(), "err.Error() GetChapterInfo controllers/readcontrollers/readcontrollers.go:32")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	response.Title = chapter.Name
	response.MangaId = chapter.Manga.Hex()
	response.Pages = chapter.Images

	var manga model.Manga
	err = manga.GetItemFromObjectId(chapter.Manga)
	if err != nil {
		log.Println(err.Error(), "err.Error() GetChapterInfo controllers/readcontrollers/readcontrollers.go:43")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	response.MangaTitle = manga.Name

	c.IndentedJSON(http.StatusOK, response)
}

func GetChapterList(c *gin.Context) {
	id := c.Param("chapterid")
	log.Println("ChapterId:", id)
	//TODO: check authorization
	fmt.Println(c.GetHeader("Authorization"))
	// fmt.Println(c.Request.Header["Authorization"])

	var response []model.Chapter

	var chapter model.Chapter
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err.Error(), "err.Error() GetChapterInfo controllers/readcontrollers/readcontrollers.go:26")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	err = chapter.GetItemFromObjectId(objId)
	if err != nil {
		log.Println(err.Error(), "err.Error() GetChapterInfo controllers/readcontrollers/readcontrollers.go:43")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	response, err = chapter.GetMangaChapterList(chapter.Manga)
	if err != nil {
		log.Println(err.Error(), "err.Error() GetChapterInfo controllers/readcontrollers/readcontrollers.go:51")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}
