package model

import (
	"GinBlog/common/errmsg"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username  string `gorm:"type:varchar(20);not null " json:"username" validate:"required,min=3,max=12" label:"用户名"`
	Password  string `gorm:"size:255;not null" json:"password" validate:"required,min=6,max=120" label:"密码"`
	Telephone string `gorm:"type:varchar(11);not null;unique" json:"telephone" validate:"required,min=6,max=120" label:"电话"`
	Role      int    `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=0" label:"角色码"`
	ImgUrl	  string `gorm:"size:255;" json:"img_url" label:"头像地址"`
}

// CheckLoginFront 前台登录
func CheckLoginFront(username string, password string) (User, int) {
	var user User
	var PasswordErr error

	db.Where("username = ?", username).First(&user)

	PasswordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if user.ID == 0 {
		return user, errmsg.ERROR_USER_NOT_EXIST
	}
	if PasswordErr != nil {
		return user, errmsg.ERROR_PASSWORD_WRONG
	}
	return user, errmsg.SUCCSE
}

// 查询手机号是否存在
func IsTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

// GetUser 查询用户
func GetUser(id int) (User, int) {
	var user User
	err := db.Limit(1).Where("ID = ?", id).Find(&user).Error
	if err != nil {
		return user, errmsg.ERROR
	}
	return user, errmsg.SUCCSE
}

//查询用户列表
func GetUsers(pageSize int, pageNum int) ([]User,int) {
	var user []User
	var total int
	err = db.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&user).Count(&total).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil,0
	}
	return user,total
}

// EditUser 编辑用户信息
func EditUser(id int, data *User) int {
	var user User
	var maps = make(map[string]interface{})
	maps["username"] = data.Username
	maps["role"] = data.Role
	err = db.Model(&user).Where("id = ? ", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	if len(user.Username)<3{
		return errmsg.ERROR_USERNAME_WRONG
	}
	return errmsg.SUCCSE
}

// ChangePassword 修改密码   TODO 待完善
func ChangePassword(id int, data *User) int {
	//var user User
	var maps = make(map[string]interface{})
	maps["password"] = data.Password

	err = db.Select("password").Where("id = ?", id).Updates(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// DeleteUser 删除用户
func DeleteUser(id int) int {
	var user User
	err = db.Where("id = ? ", id).Delete(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// CheckUpUser 更新查询
func CheckUpUser(id int, name string) (code int) {
	var user User
	db.Select("id, username").Where("username = ?", name).First(&user)
	if user.ID == uint(id) {
		return errmsg.SUCCSE
	}
	if user.ID > 0 {
		return errmsg.ERROR_USERNAME_USED //1001
	}
	return errmsg.SUCCSE
}
