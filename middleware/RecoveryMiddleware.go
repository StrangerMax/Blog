package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err!=nil{
				c.JSON(http.StatusUnauthorized, gin.H{
					"status": 400,
					"msg":  fmt.Sprint(err),
				},
				)
			}
		}()
		c.Next()
	}
}
