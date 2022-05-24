package middleware

import (
	commom "GinBlog/common"
	"GinBlog/common/errmsg"
	"GinBlog/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

//jwt中间件完成用户认证
func AuthMiddleware() gin.HandlerFunc {
	//todO 有bug
	return func(c *gin.Context) {
		//获取authorization header
		tokenString := c.GetHeader("Authorization")
		//tokenString := ctx.Request.Header.Get("x-token")
		//验证token格式 是否以Bearer开头
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			code := errmsg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": code,
				"msg":    errmsg.GetErrMsg(code)},
			)
			c.Abort() //抛弃
			return
		}
		//截取内容 Bearer占前6位
		tokenString = tokenString[7:]

		token, cliams, err := commom.ParseToken(tokenString)
		if err == errmsg.ERROR || !token.Valid {
			code := errmsg.ERROR_USER_NO_RIGHT
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": code,
				"msg":    errmsg.GetErrMsg(code)},
			)
			c.Abort() //抛弃
			return
		}
		//if time.Now().Unix() > cliams.ExpiresAt {
		//	code := errmsg.ERROR_TOKEN_RUNTIME
		//	c.JSON(http.StatusUnauthorized, gin.H{
		//		"status": code,
		//		"msg":    errmsg.GetErrMsg(code)},
		//	)
		//	c.Abort()
		//	return
		//}

		//验证通过后获取claim中的userId
		userId := cliams.UserId
		DB := model.GetDB()
		var user model.User
		DB.First(&user, userId)

		//用户不存在
		if user.ID == 0 {
			code := errmsg.ERROR_USER_NO_RIGHT
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": code,
				"msg":    errmsg.GetErrMsg(code)},
			)
			c.Abort()
			return
		}

		//用户存在  将user 信息写入上下文
		c.Set("user", user)
		c.Next()
	}
}
