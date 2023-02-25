package middleware

import (
	"errors"
	"log"
	"magna/services/chapterservice"
	"magna/services/sessionservice"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionkey := c.GetHeader("Authorization")

		if strings.Contains(c.Request.URL.String(), "/read") {
			chapterId := c.Param("chapterid")

			var userId string
			var err error
			if sessionkey == "" {
				userId = ""
			} else {
				userId, err = sessionservice.ExtractUserIdFromSessionKey(sessionkey)
				if err != nil {
					c.AbortWithError(http.StatusUnauthorized, errors.New("UnAuthorized"))
				}
			}
			owned, err := chapterservice.CheckIsOwner(chapterId, userId)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, errors.New("Error in system"))
			}
			if !owned {
				c.AbortWithError(http.StatusUnauthorized, errors.New("UnAuthorized"))
			}
			log.Println("Owned detected")

			return
		}

		valid, err := sessionservice.CheckSession(sessionkey)
		if !valid || err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("UnAuthorized"))
		}
	}
}
