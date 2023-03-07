package service

import (
	"context"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"tiktok-demo/common"
	"tiktok-demo/dao/redis"
	"tiktok-demo/logger"
	"tiktok-demo/middleware/kafka"
)

// SubscribeUser 关注用户
func SubscribeUser(userId int, toUserId int) error {
	// 1、将消息传给kafka
	// kafka关注消息格式 key: userId
	// value: toUserId:del，删除，value: toUserId:add，添加
	err = kafka.RelationClient.SendMessage(strconv.Itoa(userId), strconv.Itoa(toUserId)+":add")
	if err != nil {
		logger.Log.Error("RelationClient.SendMessage failed", zap.Any("error", err))
		return err
	}
	return nil
}

// UnsubscribeUser 取关用户
func UnsubscribeUser(userId int, toUserId int) error {
	// 1、将消息传给kafka
	// kafka关注消息格式 key: userId
	// value: toUserId:del，删除，value: toUserId:add，添加
	err = kafka.RelationClient.SendMessage(strconv.Itoa(userId), strconv.Itoa(toUserId)+":del")
	if err != nil {
		logger.Log.Error("RelationClient.SendMessage failed", zap.Any("error", err))
		return err
	}
	return nil
}

// GetFollowerList 获取粉丝列表
func GetFollowerList(userId int64) ([]common.User, error) {
	idList, err := redis.GetRelationFollowerList(context.Background(), int(userId))
	// idList, err := mysql.GetFollowedIdList(userId)

	if nil != err {
		logger.Log.Error("redis.GetRelationFollowerList failed", zap.Any("error", err))
		return nil, err
	}

	n := len(idList)
	followedList := make([]common.User, n)

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			if user, err := GetCommonUserInfoById(userId, int64(idList[i])); nil == err {
				followedList[i] = user
			} else {
				logger.Log.Error("获取用户信息失败")
			}
		}(i)
	}

	wg.Wait()
	return followedList, nil
}

// GetFollowList 获取关注列表
func GetFollowList(userId int64) ([]common.User, error) {
	idList, err := redis.GetRelationFollowList(context.Background(), int(userId))
	// idList, err := mysql.GetFollowIdList(userId)

	if nil != err {
		logger.Log.Error("redis.GetRelationFollowList failed", zap.Any("error", err))
		return nil, err
	}

	n := len(idList)
	followList := make([]common.User, n)

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			if user, err := GetCommonUserInfoById(userId, int64(idList[i])); nil == err {
				followList[i] = user
			} else {
				logger.Log.Error("获取用户信息失败")
			}
		}(i)
	}

	wg.Wait()
	return followList, nil
}
