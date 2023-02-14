package service

import (
	"errors"
	"go.uber.org/zap"
	"strconv"
	"tiktok-demo/common"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/logger"
	"tiktok-demo/middleware/kafka"
)

func AddChatRecord(userId int, toUserIdRaw, actionTypeRaw, content string) error {
	toUserId, err := strconv.Atoi(toUserIdRaw)
	if err != nil {
		return errors.New("query 'toUserId' or 'actionType' should be a integer")
	}
	actionType, err := strconv.Atoi(actionTypeRaw)
	if err != nil {
		return errors.New("query 'toUserId' or 'actionType' should be a integer")
	}
	// 对方账号是否存在
	_, err = mysql.CheckUserExist(toUserId)
	if err != nil {
		return err
	}
	// 将操作信息传给kafka，
	// key:     toUserId，
	// value:   userId:actionType:content
	err = kafka.ChatClient.SendMessage(strconv.Itoa(toUserId), strconv.Itoa(userId)+":"+strconv.Itoa(actionType)+":"+content)
	if err != nil {
		logger.Log.Error("FavoriteClient.SendMessage failed", zap.Any("error", err))
		return err
	}
	return err
}

func GetChatRecordList(userId int, toUserIdRaw string) (messageList []common.Message, err error) {
	toUserId, err := strconv.Atoi(toUserIdRaw)
	if err != nil {
		return []common.Message{}, errors.New("query 'toUserId' should be a integer")
	}
	// 对方账号是否存在
	_, err = mysql.CheckUserExist(toUserId)
	if err != nil {
		return []common.Message{}, err
	}
	messageList, err = mysql.GetChatRecordList(userId, toUserId)

	return messageList, err
}

func GetChatUnreadList(userId int, toUserIdRaw string) (messageList []common.Message, err error) {
	toUserId, err := strconv.Atoi(toUserIdRaw)
	if err != nil {
		return []common.Message{}, errors.New("query 'toUserId' should be a integer")
	}
	// 对方账号是否存在
	_, err = mysql.CheckUserExist(toUserId)
	if err != nil {
		return []common.Message{}, err
	}
	messageList, err = mysql.GetChatUnreadList(userId, toUserId)

	return messageList, err
}
