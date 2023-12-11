package routers

import (
	"magna/controllers/admincontroller"
	"magna/controllers/authcontroller"
	"magna/controllers/commentcontroller"
	"magna/controllers/favoritecontroller"
	"magna/controllers/homecontroller"
	"magna/controllers/mangacontroller"
	"magna/controllers/paymentcontroller"
	"magna/controllers/readcontroller"
	"magna/controllers/reportcontroller"
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
		apiv1Home.GET("/new/:position/:count/", homecontroller.GetNewestList)
		apiv1Home.GET("/recommend/:count/", homecontroller.GetListRecommendation)
		apiv1Home.GET("/hot/:count/", homecontroller.GetListHotItems)
		apiv1Home.GET("/user/recommend/:count/", usercontroller.GetUserRecommendation)
	}

	apiv1Manga := apiv1.Group("/manga")
	{
		apiv1Manga.GET("/:mangaid", mangacontroller.GetMangaInfo)
		apiv1Manga.GET("/:mangaid/chapterlist/:position/:count/", mangacontroller.GetMangaChapterList)
		apiv1Manga.GET("/:mangaid/commentlist/:position/:count/", mangacontroller.GetCommentList)
		apiv1Manga.POST("/:mangaid/rate", mangacontroller.SetRating)
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
		apiv1User.GET("/report", usercontroller.GetUserReport)
		apiv1User.GET("/favorite", usercontroller.GetFavoriteMangaList)
	}

	apiv1Pay := apiv1.Group("/pay")
	apiv1MomoPay := apiv1Pay.Group("/momo")
	{
		apiv1MomoPay.POST("/payurl/:chapterid", paymentcontroller.GetMomoPayURLForChapter)
		apiv1MomoPay.POST("/ipn", paymentcontroller.SetOwned)
	}

	apiv1Admin := apiv1.Group("/admin")
	apiv1User.Use(middleware.AuthMiddleware())
	{
		apiv1Admin.GET("/users", admincontroller.GetUserList)
		apiv1Admin.GET("/users/count", usercontroller.GetTotalCount)
		apiv1Admin.GET("/users/:id", admincontroller.GetUser)
		apiv1Admin.GET("/users/reference/:id", admincontroller.GetUserReference)
		apiv1Admin.POST("/users", admincontroller.CreateNewUser)
		apiv1Admin.PUT("/users/:id", admincontroller.UpdateUser)
		apiv1Admin.DELETE("/users/:id", admincontroller.DeleteUser)

		apiv1Admin.GET("/mangas", admincontroller.GetMangaList)
		apiv1Admin.GET("/mangas/count", mangacontroller.GetMangaTotalCount)
		apiv1Admin.GET("/mangas/:id", admincontroller.GetManga)
		apiv1Admin.PUT("/mangas/:id", admincontroller.UpdateManga)
		apiv1Admin.GET("/mangas/reference/:id", admincontroller.GetMangaReference)
		apiv1Admin.POST("/mangas", admincontroller.CreateNewManga)
		apiv1Admin.DELETE("/mangas/:id", admincontroller.DeleteManga)

		apiv1Admin.GET("/chapters", admincontroller.GetChapterList)
		apiv1Admin.GET("/chapters/count", mangacontroller.GetChapterTotalCount)
		apiv1Admin.GET("/chapters/:id", admincontroller.GetChapter)
		apiv1Admin.PUT("/chapters/:id", admincontroller.UpdateChapter)
		apiv1Admin.GET("/chapters/reference/:id", admincontroller.GetChapterReference)
		apiv1Admin.POST("/chapters", admincontroller.CreateNewChapter)
		apiv1Admin.DELETE("/chapters/:id", admincontroller.DeleteChapter)

		apiv1Admin.GET("/comments", admincontroller.GetCommentList)
		apiv1Admin.GET("/comments/count", commentcontroller.GetTotalCount)
		apiv1Admin.GET("/comments/:id", admincontroller.GetComment)
		apiv1Admin.GET("/comments/reference/:id", admincontroller.GetCommentReference)
		apiv1Admin.DELETE("/comments/:id", admincontroller.DeleteComment)

		apiv1Admin.GET("/reports", admincontroller.GetReportList)
		apiv1Admin.GET("/reports/count", reportcontroller.GetTotalCount)
		apiv1Admin.GET("/reports/:id", admincontroller.GetReport)
		apiv1Admin.PUT("/reports/:id", admincontroller.RespondReport)
		apiv1Admin.GET("/reports/reference/:id", admincontroller.GetReportReference)
		apiv1Admin.DELETE("/reports/:id", admincontroller.DeleteReport)

		apiv1Admin.GET("/auth", authcontroller.CheckAdmin)
	}

	apiv1.POST("/report/:chapterid", reportcontroller.CreateNewReport)
	apiv1.POST("/comment/:mangaid", commentcontroller.CreateNewComment)
	apiv1.POST("/favorite/:mangaid", favoritecontroller.SetFavorite)

	return router
}
