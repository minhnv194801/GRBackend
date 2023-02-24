package routers

import (
	"magna/controllers/homecontroller"
	"magna/controllers/mangacontroller"
	"magna/controllers/readcontroller"
	"magna/controllers/usercontroller"
	"magna/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	apiv1 := router.Group("/api/v1")

	apiv1User := apiv1.Group("/user")
	{
		apiv1User.POST("/login", usercontroller.Login)
		apiv1User.POST("/register", usercontroller.Register)
		apiv1User.GET("/refresh", usercontroller.RefreshSession)
	}

	apiv1Home := apiv1.Group("/home")
	{
		apiv1Home.GET("/count", homecontroller.GetTotalCount)
		apiv1Home.POST("/new", homecontroller.GetNewestList)
		apiv1Home.POST("/recommend", homecontroller.GetListRecommendation)
		apiv1Home.POST("/hot", homecontroller.GetListHotItems)
	}

	apiv1Manga := apiv1.Group("/manga")
	{
		apiv1Manga.GET("/:mangaid", mangacontroller.GetMangaInfo)
		apiv1Manga.POST("/:mangaid/chapterlist", mangacontroller.GetMangaChapterList)
		apiv1Manga.POST("/:mangaid/commentlist", mangacontroller.GetCommentList)
		apiv1Manga.GET("/:mangaid/comment/count", mangacontroller.GetMangaCommentCount)
	}

	apiv1Read := apiv1.Group("/read")
	{
		apiv1Read.GET("/:chapterid", readcontroller.GetChapterInfo)
		apiv1Read.GET("/:chapterid/chapterlist", readcontroller.GetChapterList)
	}

	return router
}
