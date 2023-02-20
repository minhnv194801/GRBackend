package routers

import (
	"magna/controllers/readcontrollers"
	"magna/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// corsConfig := cors.DefaultConfig()
	// corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	// corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	// corsConfig.AllowCredentials = true
	// corsConfig.AddAllowMethods("OPTIONS")

	router.Use(middleware.CORSMiddleware())

	apiv1 := router.Group("/api/v1")
	apiv1Read := apiv1.Group("/read")
	{
		apiv1Read.GET("/:chapterid", readcontrollers.GetChapterInfo)
	}

	return router
}
