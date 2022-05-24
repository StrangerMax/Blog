package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username  string `gorm:"type:varchar(20);not null " json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password  string `gorm:"size:255;not null" json:"password" validate:"required,min=6,max=120" label:"密码"`
	Telephone string `gorm:"type:varchar(11);not null;unique" json:"telephone" validate:"required,min=6,max=120" label:"电话"`
	Role      int    `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=0" label:"角色码"`
	ImgUrl	  string `gorm:"size:255;DEFAULT:https://static.nowcoder.com/head/2photo.jpg" json:"img_url" label:"头像地址"`
}
