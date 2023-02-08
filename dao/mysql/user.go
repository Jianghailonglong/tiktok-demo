package mysql

import (
	"errors"
	"gorm.io/gorm"
	"tiktok-demo/logger"
)

type User struct {
	Id       int    `gorm:"column:id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

func (User) TableName() string {
	return "users"
}

// GetUserByUserName 根据用户名查询用户
func GetUserByUserName(username string) (user User, err error) {
	res := db.Where("username = ?", username).Take(&user)
	if res.Error == gorm.ErrRecordNotFound {
		return user, errors.New("查不到该用户")
	} else {
		if res.Error != nil {
			return user, errors.New("GetUserByUserName查询失败")
		}
	}
	return
}

// InsertUser 插入新用户
func InsertUser(username, encrytedPassword string) (user User, err error) {
	user = User{
		Username: username,
		Password: encrytedPassword,
	}
	res := db.Create(&user)
	if res.Error != nil {
		return User{}, errors.New("InsertUser插入失败")
	}
	return
}

// GetUserByUserID 根据用户ID查询用户
func GetUserByUserID(userID int64) (user User, err error) {
	res := db.Where("id = ?", userID).Take(&user)
	if res.Error == gorm.ErrRecordNotFound {
		return user, errors.New("查不到该用户")
	} else {
		if res.Error != nil {
			return user, errors.New("GetUserByUserID查询失败")
		}
	}
	return
}

// GetUserByUserIDList 根据用户ID列表查询用户
func GetUserByUserIDList(userIDList []int64) (userList []User, err error) {
	userList = make([]User, len(userIDList))
	for i := 0; i < len(userIDList); i++ {
		user, err := GetUserByUserID(userIDList[i])
		if err != nil {
			logger.Log.Error("GetUserByUserID failed")
			userList[i] = user
			continue
		}
		userList[i] = user
	}
	return
}
