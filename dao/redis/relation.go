package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
	"tiktok-demo/logger"
)

// AddRelation 关注
func AddRelation(ctx context.Context, userId int, toUserId int) error {
	luaScript := redis.NewScript(`
		redis.call("SAdd", KEYS[1], ARGV[1])
        redis.call("SAdd", KEYS[2], ARGV[2])
	`)
	// 执行脚本
	_, err := luaScript.Run(ctx, client, []string{
		KeyFollowUserIdPrefix + strconv.Itoa(userId),
		KeyFollowerUserIdPrefix + strconv.Itoa(toUserId),
	}, toUserId, userId).Result()

	if err != redis.Nil {
		logger.Log.Error("client.SAdd failed", zap.Any("error", err))
		return err
	}
	logger.Log.Info("SAdd", zap.Any("key1", KeyFollowUserIdPrefix+strconv.Itoa(userId)), zap.Any("value1", toUserId), zap.Any("key2",
		KeyFollowerUserIdPrefix+strconv.Itoa(toUserId)), zap.Any("value2", userId))
	return nil
}

// CancelRelation 取消关注
func CancelRelation(ctx context.Context, userId int, toUserId int) error {
	luaScript := redis.NewScript(`
		redis.call("SRem", KEYS[1], ARGV[1])
        redis.call("SRem", KEYS[2], ARGV[2])
	`)
	// 执行脚本
	_, err := luaScript.Run(ctx, client, []string{
		KeyFollowUserIdPrefix + strconv.Itoa(userId),
		KeyFollowerUserIdPrefix + strconv.Itoa(toUserId),
	}, toUserId, userId).Result()

	if err != redis.Nil {
		logger.Log.Error("client.Rem failed", zap.Any("error", err))
		return err
	}
	logger.Log.Info("SRem", zap.Any("key1", KeyFollowUserIdPrefix+strconv.Itoa(userId)), zap.Any("value1", toUserId), zap.Any("key2",
		KeyFollowerUserIdPrefix+strconv.Itoa(toUserId)), zap.Any("value2", userId))
	return nil
}

// GetRelationFollowList 获取用户关注列表id
func GetRelationFollowList(ctx context.Context, userId int) (userIdList []int, err error) {
	key := KeyFollowUserIdPrefix + strconv.Itoa(userId)
	userIdListStr, err := client.SMembers(ctx, key).Result()
	if err != nil {
		logger.Log.Error("client.SMembers(ctx, key) failed", zap.Any("error", err))
		return nil, err
	}
	logger.Log.Info("SMembers", zap.Any("key", key), zap.Any("members", userIdListStr))
	// 转换成int
	userIdList = make([]int, len(userIdListStr))
	for i, userIdStr := range userIdListStr {
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			logger.Log.Error("strconv.Atoi(userIdStr) failed", zap.Any("error", err))
			return nil, err
		}
		userIdList[i] = userId
	}
	return userIdList, nil
}

// GetRelationFollowerList 获取用户粉丝列表id
func GetRelationFollowerList(ctx context.Context, userId int) (userIdList []int, err error) {
	key := KeyFollowerUserIdPrefix + strconv.Itoa(userId)
	userIdListStr, err := client.SMembers(ctx, key).Result()
	if err != nil {
		logger.Log.Error("client.SMembers(ctx, key) failed", zap.Any("error", err))
		return nil, err
	}
	logger.Log.Info("SMembers", zap.Any("key", key), zap.Any("members", userIdListStr))
	// 转换成int
	userIdList = make([]int, len(userIdListStr))
	for i, userIdStr := range userIdListStr {
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			logger.Log.Error("strconv.Atoi(userIdStr) failed", zap.Any("error", err))
			return nil, err
		}
		userIdList[i] = userId
	}
	return userIdList, nil
}

// CheckIsRelation 获取用户是否有关系
func CheckIsRelation(ctx context.Context, userId int, toUserId int) (bool, error) {
	key := KeyFollowUserIdPrefix + strconv.Itoa(userId)
	ok, err := client.SIsMember(ctx, key, toUserId).Result()
	if err != nil {
		logger.Log.Error("client.SMembers(ctx, key, toUserId) failed", zap.Any("error", err))
		return false, err
	}
	logger.Log.Info("SIsMember", zap.Any("key", key), zap.Any("member", toUserId), zap.Any("true", ok))
	return ok, nil
}
