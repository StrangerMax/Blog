package v1

import (
	commom "GinBlog/common"
	"GinBlog/common/errmsg"
	"GinBlog/common/validator"
	"GinBlog/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
)

var code int

//注册用户
func Register(c *gin.Context) {
	//todo 添加用户
	DB := model.GetDB()
	var requestUser = model.User{}
	c.ShouldBindJSON(&requestUser)

	//数据验证
	var msg string
	msg, code = validator.Validate(&requestUser)
	if code != errmsg.SUCCSE {
		c.JSON(http.StatusOK, gin.H{
			"status": code,
			"msg1":   msg,
		})
		return
	}

	username := requestUser.Username
	telephone := requestUser.Telephone
	password := requestUser.Password
	role := requestUser.Role

	//数据验证
	if len(telephone) != 11 {
		code = errmsg.ERROR_TELEPHONE_WRONG
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}
	if len(password) < 6 {
		code = errmsg.ERROR_PASSWORD_WRONG
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}

	log.Println(username, telephone, password, role)

	//判断手机号是否存在
	if model.IsTelephoneExist(DB, telephone) {
		code = errmsg.ERROR_TELEPHONE_USED
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}

	//创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10) //密码不为明文
	if err != nil {
		code = errmsg.ERROR_PWDHASE_WRONG
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}
	newUser := model.User{
		Username:  username,
		Telephone: telephone,
		Password:  string(hasedPassword),
		Role:      role,
	}
	DB.Create(&newUser)

	//发放token给前端
	var token string
	token, code = commom.ReleaseToken(newUser)
	if code != errmsg.SUCCSE {
		code = errmsg.ERROR_SYSTEM_WRONG
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		log.Printf("token generate error :%v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"token":  token,
		"msg":    errmsg.GetErrMsg(code),
	})
}

//登录
// Login 后台登录
func Login(c *gin.Context) {
	DB := model.GetDB()
	//使用结构体获取请求的参数
	var requestUser = model.User{}
	var token string
	c.Bind(&requestUser)

	//获取参数
	username := requestUser.Username
	password := requestUser.Password

	//数据验证
	if len(password) < 6 {
		code = errmsg.ERROR_PASSWORD_WRONG
		log.Println("密码不能少于6位")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}
	//判断用户是否存在
	var user model.User
	DB.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		code = errmsg.ERROR_USER_NOT_EXIST
		log.Println("用户不存在")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		code = errmsg.ERROR_PASSWORD_WRONG
		log.Println("密码错误")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}

	if user.Role !=1{
		code = errmsg.ERROR_USER_NO_RIGHT
		log.Println("密码错误")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": code,
			"msg":    errmsg.GetErrMsg(code),
		})
		return
	}

	//发放token给前端
	token, code = commom.ReleaseToken(user)
	if code != errmsg.SUCCSE {
		code = errmsg.ERROR_SYSTEM_WRONG
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": code,
			"msg": errmsg.GetErrMsg(code),
		})
		return
	}
	//返回结果
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"status": errmsg.SUCCSE,
		"msg":   "登录成功",
	})
}

// LoginFront 前台登录
func LoginFront(c *gin.Context) {
	var formData model.User
	_ = c.ShouldBindJSON(&formData)
	var code int

	formData, code = model.CheckLoginFront(formData.Username, formData.Password)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    formData.Username,
		"id":      formData.ID,
		"message": errmsg.GetErrMsg(code),
	})
}

// 查询用户列表
func GetUsers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))

	if pageSize == 0 {
		pageSize = -1 //limit查询-1 直接忽略该过滤
	}
	if pageNum == 0 {
		pageNum = -1
	}
	data, total := model.GetUsers(pageSize, pageNum)
	code = errmsg.SUCCSE
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":   data,
		"total":  total,
		"msg":    errmsg.GetErrMsg(code),
	})
}

// GetUserInfo 查询单个用户
func GetUserInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var maps = make(map[string]interface{})
	data, code := model.GetUser(id)
	maps["username"] = data.Username
	maps["role"] = data.Role
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    maps,
			"total":   1,
			"message": errmsg.GetErrMsg(code),
		},
	)

}

// EditUser 编辑用户
func EditUser(c *gin.Context) {
	var data model.User
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)
	code := model.CheckUpUser(id, data.Username)
	if code == errmsg.SUCCSE {
		model.EditUser(id, &data)
	}

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// ChangeUserPassword 修改密码
func ChangeUserPassword(c *gin.Context) {
	var data model.User
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)

	code = model.ChangePassword(id, &data)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	code := model.DeleteUser(id)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}
