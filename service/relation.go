package service

import (
	"tiktok-demo/dao/mysql"
)

// 关注用户
func SubscribeUser(userId int, toUserId int) (bool, error) {
	// 先看是否原有关注关系
	relation, _ := mysql.GetRelation(userId, toUserId)

	if relation == nil {
		return mysql.AddRelation(userId, toUserId)
	} else {
		return mysql.UpdateRelation(relation, mysql.SUBSCRIBED)
	}
}

// 取关用户
func UnsubscribeUser(userId int, toUserId int) (bool, error) {
	// 先看是否原有关注关系
	relation, _ := mysql.GetRelation(userId, toUserId)

	if relation == nil {
		return true, nil
	}

	return mysql.UpdateRelation(relation, mysql.UNSUBSCRIBED)
}
