package commom

import (
	"GinBlog/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"time"
)

//数据库 链接
var db *gorm.DB
var err error

func InitDb() {
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	loc := viper.GetString("datasource.loc")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		username,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc))
	db, err = gorm.Open(driverName, args)
	if err != nil {
		log.Println("连接数据库失败，请检查参数：", err)
		panic("failed to connect database,err:" + err.Error())
	}

	//gin默认表明的复数形式
	db.SingularTable(true)

	db.AutoMigrate(&model.User{})

	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	db.DB().SetMaxIdleConns(10)

	// SetMaxOpenCons 设置数据库的最大连接数量。
	db.DB().SetMaxOpenConns(100)

	// SetConnMaxLifetiment 设置连接的最大可复用时间
	db.DB().SetConnMaxLifetime(10 * time.Second)
}


func GetDB() *gorm.DB {
	return db
}
