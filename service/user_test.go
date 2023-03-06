package service

import (
	"github.com/agiledragon/gomonkey"
	"testing"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/middleware/jwt"
)

func Test_Register_When_User_Is_Not_Existed(t *testing.T) {
	gomonkey.ApplyFunc(mysql.GetUserByUserName, func(username string) (*mysql.User, error) {
		return nil, nil
	})
	gomonkey.ApplyFunc(hashAndSalt, func(password string) (string, error) {
		return "abc", nil
	})
	gomonkey.ApplyFunc(mysql.InsertUser, func(username, encrytedPassword string) (user mysql.User, err error) {
		return mysql.User{
			Password: "abc",
		}, nil
	})
	user, err := insertUser("test", "abc-raw")
	if err != nil {
		t.Errorf("want err:%v, actually get err:%v", nil, err)
	}
	if user.Password != "abc" {
		t.Errorf("want password:%v, actually get password:%v", "abc", user.Password)
	}
}

func Test_Register_When_User_Is_Existed(t *testing.T) {
	gomonkey.ApplyFunc(mysql.GetUserByUserName, func(username string) (*mysql.User, error) {
		return &mysql.User{}, nil
	})
	gomonkey.ApplyFunc(hashAndSalt, func(password string) (string, error) {
		return "abc", nil
	})
	user, token, err := Register("test", "abc-raw")
	if err != UserExistError {
		t.Errorf("want err:%v, actually get err:%v", UserExistError, err)
	}
	if token != "" {
		t.Errorf("want token:%v, actually get token:%v", "", token)
	}
	if user != (mysql.User{}) {
		t.Errorf("want user:%v, actually get user:%v", mysql.User{}, user)
	}
}

func Test_Login_When_User_Is_True_And_Password_is_True(t *testing.T) {
	gomonkey.ApplyFunc(mysql.GetUserByUserName, func(username string) (*mysql.User, error) {
		return &mysql.User{
			Username: "test",
			Password: "abc",
		}, nil
	})
	gomonkey.ApplyFunc(comparePassword, func(password, encryptedPassword string) bool {
		return true
	})
	gomonkey.ApplyFunc(jwt.GenToken, func(user mysql.User) (string, error) {
		return "", nil
	})
	user, _, err := Login("test", "abc")
	if err != nil {
		t.Errorf("want err:%v, actually get err:%v", nil, err)
	}
	if user.Username != "test" || user.Password != "abc" {
		t.Errorf("want user:%v, actually get user:%v", &mysql.User{Username: "test", Password: "abc"}, user)
	}
}

func Test_Login_When_User_Is_True_And_Password_is_False(t *testing.T) {
	gomonkey.ApplyFunc(mysql.GetUserByUserName, func(username string) (*mysql.User, error) {
		return &mysql.User{
			Username: "test",
			Password: "abc",
		}, nil
	})
	gomonkey.ApplyFunc(comparePassword, func(password, encryptedPassword string) bool {
		return false
	})
	gomonkey.ApplyFunc(jwt.GenToken, func(user mysql.User) (string, error) {
		return "", nil
	})
	user, _, err := Login("test", "abc-raw")
	if err != UserLoginDataError {
		t.Errorf("want err:%v, actually get err:%v", UserLoginDataError, err)
	}
	if user != (mysql.User{}) {
		t.Errorf("want user:%v, actually get user:%v", mysql.User{}, user)
	}
}
