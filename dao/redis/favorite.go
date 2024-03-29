package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
	"tiktok-demo/logger"
)

// AddFavorite 点赞
func AddFavorite(ctx context.Context, userId int, videoId []interface{}) error {
	key := KeyFavoriteUserIdPrefix + strconv.Itoa(userId)
	_, err := client.SAdd(ctx, key, videoId...).Result()
	if err != nil {
		logger.Log.Error("client.SAdd(ctx, key, videoId) failed", zap.Any("error", err))
		return err
	}
	logger.Log.Info("SAdd", zap.Any("key", key), zap.Any("value", videoId))
	return nil
}

// CancelFavorite 取消点赞
func CancelFavorite(ctx context.Context, userId int, videoId int) error {
	key := KeyFavoriteUserIdPrefix + strconv.Itoa(userId)
	_, err := client.SRem(ctx, key, videoId).Result()
	if err != nil {
		logger.Log.Error("client.SRem(ctx, key, videoId) failed", zap.Any("error", err))
		return err
	}
	logger.Log.Info("SRem", zap.Any("key", key), zap.Any("member", videoId))
	return nil
}

// GetUserFavoriteVideoCnt 获取用户喜欢视频的个数
func GetUserFavoriteVideoCnt(ctx context.Context, userId int) (int, error) {
	key := KeyFavoriteUserIdPrefix + strconv.Itoa(userId)
	cnt, err := client.SCard(ctx, key).Result()
	if err != nil {
		logger.Log.Error("client.SCard(ctx, key) failed", zap.Any("error", err))
		return 0, err
	}
	logger.Log.Info("SCard", zap.Any("key", key), zap.Any("cnt", cnt))
	return int(cnt), nil
}

// GetFavoriteVideoList 获取用户喜欢视频列表id
func GetFavoriteVideoList(ctx context.Context, userId int) (videoIdList []int, err error) {
	key := KeyFavoriteUserIdPrefix + strconv.Itoa(userId)
	videoIdListStr, err := client.SMembers(ctx, key).Result()
	if err != nil {
		logger.Log.Error("client.SMembers(ctx, key) failed", zap.Any("error", err))
		return nil, err
	}
	logger.Log.Info("SMembers", zap.Any("key", key), zap.Any("members", videoIdListStr))
	// 转换成int
	videoIdList = make([]int, len(videoIdListStr))
	for i, videoIdStr := range videoIdListStr {
		videoId, err := strconv.Atoi(videoIdStr)
		if err != nil {
			logger.Log.Error("strconv.Atoi(videoIdStr) failed", zap.Any("error", err))
			return nil, err
		}
		videoIdList[i] = videoId
	}
	return videoIdList, nil
}

// CheckIsFavorite 获取用户和视频是否为点赞关系
func CheckIsFavorite(ctx context.Context, userId int, videoId int) (bool, error) {
	key := KeyFavoriteUserIdPrefix + strconv.Itoa(userId)
	ok, err := client.SIsMember(ctx, key, videoId).Result()
	if err != nil {
		logger.Log.Error("client.SMembers(ctx, key) failed", zap.Any("error", err))
		return false, err
	}
	logger.Log.Info("SIsMember", zap.Any("key", key), zap.Any("member", videoId), zap.Any("true", ok))
	return ok, nil
}

// GetVideoFavoriteCount 获取视频点赞数量
func GetVideoFavoriteCount(ctx context.Context, videoId int) (int, error) {
	key := KeyFavoriteVideoIdPrefix + strconv.Itoa(videoId)
	cntStr, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		logger.Log.Error("client.Get(ctx, key) failed", zap.Any("error", err))
		return 0, err
	}
	cnt, err := strconv.Atoi(cntStr)
	if err != nil {
		logger.Log.Error("strconv.Atoi(cntStr) failed", zap.Any("error", err))
		return 0, err
	}
	logger.Log.Info("Get", zap.Any("key", key), zap.Any("value", cntStr))
	return cnt, err
}

// SetVideoFavoriteCount 设置视频点赞数量
func SetVideoFavoriteCount(ctx context.Context, videoId, cnt int) error {
	key := KeyFavoriteVideoIdPrefix + strconv.Itoa(videoId)
	_, err := client.Set(ctx, key, cnt, -1).Result()
	if err != nil {
		logger.Log.Error("client.Set(ctx, key, cnt, -1) failed", zap.Any("error", err))
		return err
	}
	logger.Log.Info("Set", zap.Any("key", key), zap.Any("value", cnt))
	return nil
}

// IncrVideoFavoriteCount 视频点赞数量+1
func IncrVideoFavoriteCount(ctx context.Context, videoId int) error {
	key := KeyFavoriteVideoIdPrefix + strconv.Itoa(videoId)
	_, err := client.Incr(ctx, key).Result()
	if err != nil {
		logger.Log.Error("client.Incr(ctx, key) failed", zap.Any("error", err))
		return err
	}
	logger.Log.Info("Incr", zap.Any("key", key))
	return nil
}

// DecrVideoFavoriteCount 视频点赞数量-1
func DecrVideoFavoriteCount(ctx context.Context, videoId int) error {
	key := KeyFavoriteVideoIdPrefix + strconv.Itoa(videoId)
	_, err := client.Decr(ctx, key).Result()
	if err != nil {
		logger.Log.Error("client.Decr(ctx, key) failed", zap.Any("error", err))
		return err
	}
	logger.Log.Info("Decr", zap.Any("key", key))
	return nil
}

// DelVideoFavoriteCount 删除视频点赞数量
func DelVideoFavoriteCount(ctx context.Context, videoId int) error {
	key := KeyFavoriteVideoIdPrefix + strconv.Itoa(videoId)
	_, err := client.Del(ctx, key).Result()
	if err != nil {
		logger.Log.Error("client.Del(ctx, key) failed", zap.Any("error", err))
		return err
	}
	logger.Log.Info("Del", zap.Any("key", key))
	return nil
}
