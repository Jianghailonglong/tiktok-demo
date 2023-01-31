package service

import "errors"

var (
	UserExistError     = errors.New("用户名重复")
	UserLoginDataError = errors.New("用户或密码错误")
	InsertDataError    = errors.New("插入数据错误")
	UpdateDataError    = errors.New("更新数据错误")
	DeleteDataError    = errors.New("删除数据错误")
	QueryDataError     = errors.New("查询数据错误")
)
