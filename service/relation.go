package service

import (
	"sync"
	"tiktok-demo/common"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/logger"
)

// SubscribeUser 关注用户
func SubscribeUser(userId int, toUserId int) error {
	// 先看是否原有关注关系
	relation, _ := mysql.GetRelation(userId, toUserId)

	if relation == nil {
		return mysql.AddRelation(userId, toUserId)
	} else {
		return mysql.UpdateRelation(relation, mysql.SUBSCRIBED)
	}
}

// UnsubscribeUser 取关用户
func UnsubscribeUser(userId int, toUserId int) error {
	// 先看是否原有关注关系
	relation, _ := mysql.GetRelation(userId, toUserId)

	if relation == nil {
		return nil
	}

	return mysql.UpdateRelation(relation, mysql.UNSUBSCRIBED)
}

// GetFollowerList 获取粉丝列表
func GetFollowerList(userId int64) ([]common.User, error) {
	idList, err := mysql.GetFollowedIdList(userId)

	if nil != err {
		return nil, err
	}

	n := len(idList)
	followedList := make([]common.User, n)

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			if user, err := GetCommonUserInfoById(userId, idList[i]); nil == err {
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
	idList, err := mysql.GetFollowIdList(userId)

	if nil != err {
		return nil, err
	}

	n := len(idList)
	followList := make([]common.User, n)

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			if user, err := GetCommonUserInfoById(userId, idList[i]); nil == err {
				followList[i] = user
			} else {
				logger.Log.Error("获取用户信息失败")
			}
		}(i)
	}

	wg.Wait()
	return followList, nil
}
