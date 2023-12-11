package readcontroller

import (
	"log"
	"magna/responses"
	"magna/services/chapterservice"
	"magna/services/mangaservice"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetChapterInfo(c *gin.Context) {
	id := c.Param("chapterid")

	var response responses.ReadResponse

	chapter, err := chapterservice.GetChapterInfo(id)
	if err != nil {
		log.Println(err.Error(), "err.Error() GetChapterInfo controllers/readcontrollers/readcontrollers.go:26")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	response.Title = chapter.Name
	response.MangaId = chapter.Manga.Hex()
	response.Pages = chapter.Images
	response.Price = chapter.Price

	manga, err := mangaservice.GetItemFromId(chapter.Manga.Hex())
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

	response, err := chapterservice.GetSameMangaChapterList(id)
	if err != nil {
		log.Println(err.Error(), "err.Error() GetChapterInfo controllers/readcontrollers/readcontrollers.go:32")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}
