package routers

import (
	"GinBlog/controller/v1"
	"GinBlog/middleware"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func createMyRender() multitemplate.Renderer {
	p := multitemplate.NewRenderer()
	p.AddFromFiles("admin", "web/admin/dist/index.html")
	p.AddFromFiles("front", "web/front/dist/index.html")
	return p
}

func Router() {
	gin.SetMode("debug")
	r := gin.New()
	r.HTMLRender = createMyRender()
	r.Use(middleware.CORSMiddleware(), middleware.RecoveryMiddleware(), middleware.Logger())

	r.Static("/static", "./web/front/dist/static")
	r.Static("/admin", "./web/admin/dist")
	r.StaticFile("/favicon.ico", "/web/front/dist/favicon.ico")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "front", nil)
	})

	r.GET("/admin", func(c *gin.Context) {
		c.HTML(200, "admin", nil)
	})

	/*
		后台管理路由接口
	*/
	auth := r.Group("api/v1")
	//auth.Use(middleware.AuthMiddleware())
	{
		// 用户模块的路由接口
		auth.GET("admin/users", v1.GetUsers)
		auth.PUT("user/:id", v1.EditUser)      //TODO bug存在 密码 手机号等修改  必须提供name role字段
		auth.DELETE("user/:id", v1.DeleteUser) // 软删除
		////修改密码
		auth.PUT("admin/changepw/:id", v1.ChangeUserPassword) //TODO 存在bug  该方法需要对传入的数据进行加密解密
		//// 分类模块的路由接口
		auth.GET("admin/category", v1.GetCate)
		auth.POST("category/add", v1.AddCate)
		auth.PUT("category/:id", v1.EditCate)
		auth.DELETE("category/:id", v1.DeleteCate)
		//// 文章模块的路由接口
		auth.GET("admin/article/info/:id", v1.GetArtInfo) // GetArtInfo 查询单个文章信息
		auth.GET("admin/article", v1.GetArt)              // GetArt 查询文章列表   TODO bug 无法查询
		auth.POST("article/add", v1.AddArticle)           // AddArticle 添加文章
		auth.PUT("article/:id", v1.EditArt)               // EditArt 编辑文章  TODO bug 无法修改
		auth.DELETE("article/:id", v1.DeleteArt)          // DeleteArt 删除文章 TODO bug 不存在的删除没有正确提示
		//// 上传文件
		auth.POST("upload", v1.UpLoad)
		//// 更新个人设置
		auth.GET("admin/profile/:id", v1.GetProfile) // GetProfile 获取个人信息设置
		auth.PUT("profile/:id", v1.UpdateProfile)    // UpdateProfile 更新个人信息设置
		//// 评论模块
		//auth.GET("comment/list", v1.GetCommentList)
		//auth.DELETE("delcomment/:id", v1.DeleteComment)
		//auth.PUT("checkcomment/:id", v1.CheckComment)
		//auth.PUT("uncheckcomment/:id", v1.UncheckComment)
	}

	/*
		前端展示页面接口
	*/
	router := r.Group("api/v1")
	{
		// 用户信息模块
		router.POST("user/add", v1.Register)   //TODO 账户创建对于软删除的手机号存在bug
		router.GET("user/:id", v1.GetUserInfo) //TODO 账户信息查询只提供name role 差手机号
		router.GET("users", v1.GetUsers)       //查询单个用户信息
		//
		//// 文章分类信息模块
		router.GET("category", v1.GetCate)
		router.GET("category/:id", v1.GetCateInfo)
		//
		//// 文章模块
		router.GET("article", v1.GetArt)              // GetArt 查询文章列表  TODO bug 无法查询
		router.GET("article/list/:id", v1.GetCateArt) // GetCateArt 查询分类下的所有文章
		router.GET("article/info/:id", v1.GetArtInfo) // GetArtInfo 查询单个文章信息
		//
		//// 登录控制模块
		router.POST("login", v1.Login)
		router.POST("loginfront", v1.LoginFront)
		//
		//// 获取个人设置信息
		router.GET("profile/:id", v1.GetProfile) // GetProfile 获取个人信息设置
		//
		//// 评论模块
		//router.POST("addcomment", v1.AddComment)
		//router.GET("comment/info/:id", v1.GetComment)
		//router.GET("commentfront/:id", v1.GetCommentListFront)
		//router.GET("commentcount/:id", v1.GetCommentCount)
	}
	server := viper.GetString("server.port")
	r.Run(":" + server)
}
