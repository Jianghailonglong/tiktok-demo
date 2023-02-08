package service

import (
	"sync"
	"tiktok-demo/common"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/logger"
)

const (
	FOLLOW   = 1 // 获取关注列表
	FOLLOWER = 2 // 获取粉丝列表
)

// 关注用户
func SubscribeUser(userId int, toUserId int) error {
	// 先看是否原有关注关系
	relation, _ := mysql.GetRelation(userId, toUserId)

	if relation == nil {
		return mysql.AddRelation(userId, toUserId)
	} else {
		return mysql.UpdateRelation(relation, mysql.SUBSCRIBED)
	}
}

// 取关用户
func UnsubscribeUser(userId int, toUserId int) error {
	// 先看是否原有关注关系
	relation, _ := mysql.GetRelation(userId, toUserId)

	if relation == nil {
		return nil
	}

	return mysql.UpdateRelation(relation, mysql.UNSUBSCRIBED)
}

// 获取粉丝列表
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
			if user, err := GetInfoById(idList[i], userId, FOLLOWER); nil == err {

				followedList[i] = user
			}
		}(i)
	}

	wg.Wait()
	return followedList, nil
}

// 获取关注列表
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
			if user, err := GetInfoById(idList[i], userId, FOLLOW); nil == err {

				followList[i] = user
			}
		}(i)
	}

	wg.Wait()
	return followList, nil
}

// 根据id获取单个用户的所有信息, relationType= 1代表关注列表 2代表粉丝列表
func GetInfoById(userId int64, withUserId int64, relationType int) (common.User, error) {
	user, err := mysql.GetInfoById(userId, withUserId, relationType)
	if nil != err {
		logger.Log.Error(err.Error())
		return user, err
	}
	return user, nil
}
