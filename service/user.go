package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"tiktok-demo/common"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/dao/redis"
	"tiktok-demo/logger"
	"tiktok-demo/middleware/jwt"
)

// Register 注册用户
func Register(username, password string) (mysql.User, string, error) {
	// 1、判断username是否存在，不存在返回true
	if user, _ := mysql.GetUserByUserName(username); user != nil {
		return mysql.User{}, "", UserExistError
	}
	// 2、插入用户和加密的密码
	user, err := insertUser(username, password)
	if err != nil {
		logger.Log.Error(err.Error())
		return mysql.User{}, "", err
	}
	// 3、生成token
	token, err := jwt.GenToken(user)
	if err != nil {
		logger.Log.Error(err.Error())
		return mysql.User{}, "", err
	}
	return user, token, nil
}

// insertUser 插入新用户
func insertUser(username, password string) (mysql.User, error) {
	encryptedPassword, err := hashAndSalt(password)
	if err != nil {
		return mysql.User{}, err
	}
	return mysql.InsertUser(username, encryptedPassword)
}

// hashAndSalt 加密密码
func hashAndSalt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Login 登录账号
func Login(username, password string) (mysql.User, string, error) {
	// 1、查询用户
	user, err := mysql.GetUserByUserName(username)
	if user == nil {
		logger.Log.Error(err.Error())
		return mysql.User{}, "", err
	}
	// 2、比较用户和密码
	if username != user.Username || !comparePassword(password, user.Password) {
		return mysql.User{}, "", UserLoginDataError
	}
	// 3、生成token
	token, err := jwt.GenToken(*user)
	if err != nil {
		logger.Log.Error(err.Error())
		return mysql.User{}, "", err
	}
	return *user, token, nil
}

func comparePassword(password, encryptedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// GetCommonUserInfoById 获取common.User
func GetCommonUserInfoById(userId int64, withUserId int64) (common.User, error) {
	// 1、根据id获取单个用户的信息（粉丝数、关注数）
	user, err := mysql.GetInfoById(userId, withUserId)
	if nil != err {
		logger.Log.Error(err.Error())
		return user, err
	}
	// 2、获取作品数以及用户发布视频id列表
	videoIdList, err := redis.GetPublishVideoList(context.Background(), int(withUserId))
	if err != nil {
		return user, err
	}
	user.WorkCount = int64(len(videoIdList))
	// 3、获取用户点赞数，其实就是喜欢列表的数目
	favoriteCount, err := redis.GetUserFavoriteVideoCnt(context.Background(), int(withUserId))
	if err != nil {
		return user, err
	}
	user.FavoriteCount = int64(favoriteCount)
	// 4、获取用户被点赞总数目，实质上是用户发布视频被点赞总数
	totalFavorited := 0
	for _, videoId := range videoIdList {
		count, err := redis.GetVideoFavoriteCount(context.Background(), videoId)
		if err != nil {
			continue
		}
		totalFavorited += count
	}
	user.TotalFavorited = int64(totalFavorited)
	return user, nil
}
