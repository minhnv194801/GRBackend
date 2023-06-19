package middleware

import (
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, User-Agent, Origin, Authorization,Content-Length, Access-Control-Allow-Origin, content-type, Content-Type, Accept, accept, Accept-Encoding,Accept-Language, sessionkey, token, Connection, Sec-WebSocket-Extensions, Sec-WebSocket-Key, Sec-WebSocket-Version, Upgrade, Sec-Fetch-Mode, Sec-Fetch-Dest,Sec-Fetch-Site, Referer, Access-Control-Request-Method, Access-Control-Request-Headers, range")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS,GET,HEAD,DELETE,PUT")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Range")

		if c.Request.Method == "OPTIONS" {
			// for k, vals := range c.Request.Header {
			// 	log.Printf("%s", k)
			// 	for _, v := range vals {
			// 		log.Printf("\t%s", v)
			// 	}
			// }
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
