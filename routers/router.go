package routers

import (
	"magna/controllers/authcontroller"
	"magna/controllers/homecontroller"
	"magna/controllers/mangacontroller"
	"magna/controllers/readcontroller"
	"magna/controllers/searchcontroller"
	"magna/controllers/usercontroller"
	"magna/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	apiv1 := router.Group("/api/v1")

	apiv1Auth := apiv1.Group("/auth")
	{
		apiv1Auth.POST("/login", authcontroller.Login)
		apiv1Auth.POST("/register", authcontroller.Register)
		apiv1Auth.GET("/refresh", authcontroller.RefreshSession)
	}

	apiv1Home := apiv1.Group("/home")
	{
		apiv1Home.POST("/new", homecontroller.GetNewestList)
		apiv1Home.POST("/recommend", homecontroller.GetListRecommendation)
		apiv1Home.POST("/hot", homecontroller.GetListHotItems)
	}

	apiv1Manga := apiv1.Group("/manga")
	{
		apiv1Manga.GET("/:mangaid", mangacontroller.GetMangaInfo)
		apiv1Manga.POST("/:mangaid/chapterlist", mangacontroller.GetMangaChapterList)
		apiv1Manga.POST("/:mangaid/commentlist", mangacontroller.GetCommentList)
	}

	apiv1Read := apiv1.Group("/read")
	apiv1Read.Use(middleware.AuthMiddleware())
	{
		apiv1Read.GET("/:chapterid", readcontroller.GetChapterInfo)
		apiv1Read.GET("/:chapterid/chapterlist", readcontroller.GetChapterList)
	}

	apiv1.POST("/search", searchcontroller.Search)

	apiv1User := apiv1.Group("/user")
	apiv1User.Use(middleware.AuthMiddleware())
	{
		apiv1User.GET("/info", usercontroller.GetUserInfo)
		apiv1User.POST("/info", usercontroller.UpdateUserInfo)
		apiv1User.GET("/owned", usercontroller.GetOwnedChapter)
	}

	return router
}
