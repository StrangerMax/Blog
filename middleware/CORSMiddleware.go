package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func CORSMiddleware() gin.HandlerFunc {
	//return func(c *gin.Context) {
	//	c.Writer.Header().Set("Access-Control-Allow-Origin","http://localhost:8080") //同源  *表示所有域名都可以访问
	//	c.Writer.Header().Set("Access-Control-Allow-Max-Age","84600")
	//	c.Writer.Header().Set("Access-Control-Allow-Methods","*")
	//	c.Writer.Header().Set("Access-Control-Allow-Headers","*")
	//	c.Writer.Header().Set("Access-Control-Allow-Credentials","true")
	//
	//	if c.Request.Method == http.MethodOptions{
	//		c.AbortWithStatus(http.StatusOK)
	//	}else {
	//		c.Next()
	//	}
	//}
	return func(c *gin.Context) {
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"}, // 等同于允许所有域名 #AllowAllOrigins:  true
			AllowMethods:     []string{"*"}, //"GET", "POST", "PUT", "DELETE", "OPTIONS"等
			AllowHeaders:     []string{"*", "Authorization"},
			ExposeHeaders:    []string{"Content-Length", "text/plain", "Authorization", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		})
	}
}
