package service

import (
	"errors"
	"strconv"
	"tiktok-demo/common"
	"tiktok-demo/dao/mysql"
)

func AddChatRecord(userId int, toUserIdRaw, actionTypeRaw, content string) error {
	toUserId,err:=strconv.Atoi(toUserIdRaw)
	if err!=nil{
		return errors.New("query 'toUserId' or 'actionType' should be a integer")
	}
	actionType,err:=strconv.Atoi(actionTypeRaw)
	if err!=nil{
		return errors.New("query 'toUserId' or 'actionType' should be a integer")
	}
	// 对方账号是否存在
	_,err=mysql.CheckUserExist(toUserId)
	if err!=nil{
		return err
	}
	err = mysql.InsertChatRecord(userId, toUserId, actionType, content)
	return err
}

func GetChatRecordList(userId int, toUserIdRaw string)(messageList []common.Message,err error){
	toUserId,err:=strconv.Atoi(toUserIdRaw)
	if err!=nil{
		return []common.Message{},errors.New("query 'toUserId' should be a integer")
	}
	// 对方账号是否存在
	_,err=mysql.CheckUserExist(toUserId)
	if err!=nil{
		return []common.Message{},err
	}
	messageList,err=mysql.GetChatRecordList(userId,toUserId)
	return messageList,err
}