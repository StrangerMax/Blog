package commom

import (
	"GinBlog/model"
	"GinBlog/common/errmsg"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var JwtKey = []byte("GinBlog_secret")

type Claims struct {
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

//生成token
func ReleaseToken(user model.User) (string, int) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) //token过期时间 7天
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(), //token发放时间
			Issuer:    "GinBlog.tech",    //发放方
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", errmsg.ERROR
	}
	return tokenString, errmsg.SUCCSE
}

//验证token
func ParseToken(tokenString string) (*jwt.Token, *Claims, int) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return JwtKey, nil
	})
	if err != nil {
		return token, claims, errmsg.ERROR
	} else {
		return token, nil, errmsg.SUCCSE
	}
}
