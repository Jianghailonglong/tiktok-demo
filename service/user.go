package service

import (
	"golang.org/x/crypto/bcrypt"
	"tiktok-demo/common"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/logger"
	"tiktok-demo/middleware/jwt"
)

// Register 注册用户
func Register(username, password string) (mysql.User, string, error) {
	// 1、判断username是否存在，不存在返回true
	if isExisted, err := checkUserExist(username); err != nil {
		logger.Log.Error(err.Error())
		return mysql.User{}, "", err
	} else {
		if isExisted {
			return mysql.User{}, "", UserExistError
		}
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

// checkUserExist 判断用户是否存在
func checkUserExist(username string) (bool, error) {
	user, err := mysql.GetUserByUserName(username)
	if err != nil {
		return false, err
	}
	if user != (mysql.User{}) {
		return true, nil
	}
	return false, nil
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
	if err != nil {
		logger.Log.Error(err.Error())
		return mysql.User{}, "", err
	}
	// 2、比较用户和密码
	if username != user.Username || !comparePassword(password, user.Password) {
		return mysql.User{}, "", UserLoginDataError
	}
	// 3、生成token
	token, err := jwt.GenToken(user)
	if err != nil {
		logger.Log.Error(err.Error())
		return mysql.User{}, "", err
	}
	return user, token, nil
}

func comparePassword(password, encryptedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// GetUserInfo 获取userInfo
func GetUserInfo(userID int64) (common.User, error) {
	user, err := mysql.GetUserByUserID(userID)
	if err != nil {
		logger.Log.Error(err.Error())
		return common.User{}, err
	}
	userInfo := common.User{
		Id:   int64(user.Id),
		Name: user.Username,
		// TODO 关注者、粉丝等
	}
	return userInfo, nil
}
