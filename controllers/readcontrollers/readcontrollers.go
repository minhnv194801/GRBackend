package readcontrollers

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
	err = chapter.GetFromObjectId(objId)
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

	chapterList, err := chapter.GetItemListFromObjectId(manga.Chapters)
	if err != nil {
		log.Println(err.Error(), "err.Error() GetChapterInfo controllers/readcontrollers/readcontrollers.go:51")
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	response.ChapterList = chapterList

	// multiFilter := bson.M{"_id": bson.M{"$in": mangaDoc.Chapters}}
	// multiFindOpts := options.Find().SetSort(bson.D{{"updateTime", 1}})
	// cursor, err := chapterColl.Find(context.TODO(), multiFilter, multiFindOpts)
	// if err != nil {
	// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
	// 	return
	// }
	// defer cursor.Close(context.Background())
	// for cursor.Next(context.Background()) {
	// 	var temp model.Chapter
	// 	cursor.Decode(&temp)
	// 	var chapterItem responses.ReadChapterListItem
	// 	chapterItem.Id = temp.Id.Hex()
	// 	chapterItem.Title = temp.Name
	// 	response.ChapterList = append(response.ChapterList, chapterItem)
	// }

	c.IndentedJSON(http.StatusOK, response)
}
