package jwt

import (
	"tiktok-demo/conf"
	"tiktok-demo/dao/mysql"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	// ContextUserIDKey 验证token成功后，将userID加入到context上下文中
	ContextUserIDKey = "userID"
)

type MyClaims struct {
	UserId int `json:"id"`
	jwt.StandardClaims
}

// GenToken 登录、注册生成token
func GenToken(user mysql.User) (string, error) {
	claims := MyClaims{
		user.Id,
		jwt.StandardClaims{
			Audience:  user.Username,
			ExpiresAt: time.Now().Unix() + int64(conf.Config.JwtExpire),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			Subject:   "token",
			Issuer:    "tiktok-demo",
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(conf.Config.JwtSecret))
	// token格式：Bearer SDFFJFFFJ
	return "Bearer " + token, err
}

func parseToken(tokenString string) (claims *MyClaims, err error) {
	claims = &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(conf.Config.JwtSecret), nil
	})
	if err == nil && token != nil {
		if claim, ok := token.Claims.(*MyClaims); ok && token.Valid {
			return claim, nil
		}
	}
	return
}
