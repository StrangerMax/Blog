package model

import (
	"GinBlog/common/errmsg"
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/viper"
	"mime/multipart"
)

var (
	AccessKey  string
	SecretKey  string
	Bucket     string
	QiniuSever string
)

func UploadFile(file multipart.File, fileSize int64) (string, int) {
	LoadQiniu()

	putPolicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac := qbox.NewMac(AccessKey, SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}

	putExtra := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err = formUploader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	if err != nil {
		return "", errmsg.ERROR
	}
	url := QiniuSever + ret.Key
	return url, errmsg.SUCCSE
}

func LoadQiniu()  {
	InitConfig()
	AccessKey = viper.GetString("qiniu.AccessKey")
	SecretKey = viper.GetString("qiniu.SecretKey")
	Bucket = viper.GetString("qiniu.Bucket")
	QiniuSever = viper.GetString("qiniu.QiniuSever")
}
