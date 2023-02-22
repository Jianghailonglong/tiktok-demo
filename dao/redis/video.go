package redis

import (
	"context"
	"go.uber.org/zap"
	"strconv"
	"tiktok-demo/logger"
)

// AddVideo 用户发布视频
func AddVideo(ctx context.Context, userId int, videoId int) error {
	key := KeyPublishUserIdPrefix + strconv.Itoa(userId)
	_, err := client.SAdd(ctx, key, videoId).Result()
	if err != nil {
		logger.Log.Error("client.SAdd(ctx, key, videoId) failed", zap.Any("error", err))
		return err
	}
	logger.Log.Info("SAdd", zap.Any("key", key), zap.Any("videoId", videoId))
	return nil
}

// GetPublishVideoList 获取用户发布视频列表id
func GetPublishVideoList(ctx context.Context, userId int) (videoIdList []int, err error) {
	key := KeyPublishUserIdPrefix + strconv.Itoa(userId)
	videoIdListStr, err := client.SMembers(ctx, key).Result()
	if err != nil {
		logger.Log.Error("client.SMembers(ctx, key) failed", zap.Any("error", err))
		return nil, err
	}
	logger.Log.Info("SMembers", zap.Any("key", key), zap.Any("videoIdListStr", videoIdListStr))
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
