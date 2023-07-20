package admincontroller

import (
	"log"
	"magna/model"
	"magna/requests"
	"magna/services/chapterservice"
	"magna/services/mangaservice"
	"magna/services/userservice"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserList(c *gin.Context) {
	listRange := c.Request.URL.Query().Get("range")
	listRange = strings.Trim(listRange, "[")
	listRange = strings.Trim(listRange, "]")
	pos, _ := strconv.Atoi(strings.Split(listRange, ",")[0])
	count, _ := strconv.Atoi(strings.Split(listRange, ",")[1])

	listSort := c.Request.URL.Query().Get("sort")
	listSort = strings.Trim(listSort, "[")
	listSort = strings.Trim(listSort, "]")
	sortField := strings.Split(listSort, ",")[0]
	sortField = strings.Trim(sortField, "\"")
	sortType := strings.Split(listSort, ",")[1]
	sortType = strings.Trim(sortType, "\"")

	filter := c.Request.URL.Query().Get("filter")
	filter = strings.Trim(filter, "[")
	filter = strings.Trim(filter, "]")
	var filterField string
	var filterValue string
	if len(strings.Split(filter, ",")) >= 2 {
		filterField = strings.Split(filter, ",")[0]
		filterField = strings.Trim(filterField, "\"")
		filterValue = strings.Split(filter, ",")[1]
		filterValue = strings.Trim(filterValue, "\"")
	}

	var userList []model.User
	var totalCount int
	var err error

	userList, totalCount, err = new(model.User).GetItemList(pos, count, sortField, sortType)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	endIndex := pos + len(userList) - 1
	c.Header("Content-Range", strconv.Itoa(pos)+"-"+strconv.Itoa(endIndex)+"/"+strconv.Itoa(totalCount))

	c.IndentedJSON(http.StatusOK, userList)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")

	user := new(model.User)
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}
	err = user.GetItemFromObjectId(objId)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	c.IndentedJSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	err = new(model.User).DeleteUserById(objId)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func CreateNewUser(c *gin.Context) {
	req := requests.AdminCreateAccountRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal system error"})
		return
	}

	id, err := userservice.CreateAccount(req.Email, req.Password, req.Role)
	if err != nil {
		if err.Error() == "Email đã tồn tại" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"errors": gin.H{"Email": err.Error()}})
			return
		}
		if err.Error() == "Password không hợp lệ" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"errors": gin.H{"Password": err.Error()}})
			return
		}
		if err.Error() == "Email không hợp lệ" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"errors": gin.H{"Email": err.Error()}})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"id": id})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func GetMangaList(c *gin.Context) {
	listRange := c.Request.URL.Query().Get("range")
	listRange = strings.Trim(listRange, "[")
	listRange = strings.Trim(listRange, "]")
	pos, _ := strconv.Atoi(strings.Split(listRange, ",")[0])
	count, _ := strconv.Atoi(strings.Split(listRange, ",")[1])

	listSort := c.Request.URL.Query().Get("sort")
	listSort = strings.Trim(listSort, "[")
	listSort = strings.Trim(listSort, "]")
	sortField := strings.Split(listSort, ",")[0]
	sortField = strings.Trim(sortField, "\"")
	sortType := strings.Split(listSort, ",")[1]
	sortType = strings.Trim(sortType, "\"")

	filter := c.Request.URL.Query().Get("filter")
	filter = strings.Trim(filter, "[")
	filter = strings.Trim(filter, "]")
	var filterField string
	var filterValue string
	if len(strings.Split(filter, ",")) >= 2 {
		filterField = strings.Split(filter, ",")[0]
		filterField = strings.Trim(filterField, "\"")
		filterValue = strings.Split(filter, ",")[1]
		filterValue = strings.Trim(filterValue, "\"")
	}

	var mangaList []model.Manga
	var totalCount int
	var err error
	switch filterField {
	case "name":
		mangaList, totalCount, err = new(model.Manga).GetItemListFilterByName(pos, count, sortField, sortType, filterValue)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		}
		break
	case "followedUsers":
		mangaList, totalCount, err = new(model.Manga).GetItemListFilterByFollowedUsers(pos, count, sortField, sortType, filterValue)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		}
		break
	default:
		mangaList, totalCount, err = new(model.Manga).GetItemList(pos, count, sortField, sortType)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		}
	}

	endIndex := pos + len(mangaList) - 1
	c.Header("Content-Range", strconv.Itoa(pos)+"-"+strconv.Itoa(endIndex)+"/"+strconv.Itoa(totalCount))

	c.IndentedJSON(http.StatusOK, mangaList)
}

func GetManga(c *gin.Context) {
	id := c.Param("id")

	manga := new(model.Manga)
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}
	err = manga.GetItemFromObjectId(objId)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	c.IndentedJSON(http.StatusOK, manga)
}

func DeleteManga(c *gin.Context) {
	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	err = new(model.Manga).DeleteMangaById(objId)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func CreateNewManga(c *gin.Context) {
	req := requests.AdminCreateMangaRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal system error"})
		return
	}
	log.Println(req)
	id, err := mangaservice.CreateManga(req.Name, strings.Split(req.AlternateName, ","), req.Author, req.Cover, req.Description, req.IsRecommend, strings.Split(req.Tags, ","))
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal system error"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func GetChapterList(c *gin.Context) {
	listRange := c.Request.URL.Query().Get("range")
	listRange = strings.Trim(listRange, "[")
	listRange = strings.Trim(listRange, "]")
	pos, _ := strconv.Atoi(strings.Split(listRange, ",")[0])
	count, _ := strconv.Atoi(strings.Split(listRange, ",")[1])

	listSort := c.Request.URL.Query().Get("sort")
	listSort = strings.Trim(listSort, "[")
	listSort = strings.Trim(listSort, "]")
	sortField := strings.Split(listSort, ",")[0]
	sortField = strings.Trim(sortField, "\"")
	sortType := strings.Split(listSort, ",")[1]
	sortType = strings.Trim(sortType, "\"")

	chapterList, totalCount, err := new(model.Chapter).GetItemList(pos, count, sortField, sortType)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	endIndex := pos + len(chapterList) - 1
	c.Header("Content-Range", strconv.Itoa(pos)+"-"+strconv.Itoa(endIndex)+"/"+strconv.Itoa(totalCount))

	c.IndentedJSON(http.StatusOK, chapterList)
}

func GetChapter(c *gin.Context) {
	id := c.Param("id")

	chapter := new(model.Chapter)
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}
	err = chapter.GetItemFromObjectId(objId)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	c.IndentedJSON(http.StatusOK, chapter)
}

func DeleteChapter(c *gin.Context) {
	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	err = new(model.Chapter).DeleteChapterById(objId)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func CreateNewChapter(c *gin.Context) {
	req := requests.AdminCreateChapterRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal system error"})
		return
	}

	id, err := chapterservice.CreateChapter(req.MangaId, req.Title, req.Cover, req.Price, strings.Split(req.Images, ","))
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal system error"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func GetCommentList(c *gin.Context) {
	listRange := c.Request.URL.Query().Get("range")
	listRange = strings.Trim(listRange, "[")
	listRange = strings.Trim(listRange, "]")
	pos, _ := strconv.Atoi(strings.Split(listRange, ",")[0])
	count, _ := strconv.Atoi(strings.Split(listRange, ",")[1])

	listSort := c.Request.URL.Query().Get("sort")
	listSort = strings.Trim(listSort, "[")
	listSort = strings.Trim(listSort, "]")
	sortField := strings.Split(listSort, ",")[0]
	sortField = strings.Trim(sortField, "\"")
	sortType := strings.Split(listSort, ",")[1]
	sortType = strings.Trim(sortType, "\"")

	commentList, totalCount, err := new(model.Comment).GetItemList(pos, count, sortField, sortType)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	endIndex := pos + len(commentList) - 1
	c.Header("Content-Range", strconv.Itoa(pos)+"-"+strconv.Itoa(endIndex)+"/"+strconv.Itoa(totalCount))

	c.IndentedJSON(http.StatusOK, commentList)
}

func GetComment(c *gin.Context) {
	id := c.Param("id")

	comment := new(model.Comment)
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}
	err = comment.GetItemFromObjectId(objId)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	c.IndentedJSON(http.StatusOK, comment)
}

func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	err = new(model.Comment).DeleteCommentById(objId)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func GetReportList(c *gin.Context) {
	listRange := c.Request.URL.Query().Get("range")
	listRange = strings.Trim(listRange, "[")
	listRange = strings.Trim(listRange, "]")
	pos, _ := strconv.Atoi(strings.Split(listRange, ",")[0])
	count, _ := strconv.Atoi(strings.Split(listRange, ",")[1])

	listSort := c.Request.URL.Query().Get("sort")
	listSort = strings.Trim(listSort, "[")
	listSort = strings.Trim(listSort, "]")
	sortField := strings.Split(listSort, ",")[0]
	sortField = strings.Trim(sortField, "\"")
	sortType := strings.Split(listSort, ",")[1]
	sortType = strings.Trim(sortType, "\"")

	reportList, totalCount, err := new(model.Report).GetItemList(pos, count, sortField, sortType)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	endIndex := pos + len(reportList) - 1
	c.Header("Content-Range", strconv.Itoa(pos)+"-"+strconv.Itoa(endIndex)+"/"+strconv.Itoa(totalCount))

	c.IndentedJSON(http.StatusOK, reportList)
}

func GetReport(c *gin.Context) {
	id := c.Param("id")

	report := new(model.Report)
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}
	err = report.GetItemFromObjectId(objId)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	c.IndentedJSON(http.StatusOK, report)
}

func DeleteReport(c *gin.Context) {
	id := c.Param("id")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	err = new(model.Report).DeleteReportById(objId)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}
