package kafka

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"tiktok-demo/dao/redis"
	"tiktok-demo/logger"
	"time"
)

// FavoriteConsumerGroup 点赞相关的消费者
type FavoriteConsumerGroup struct {
	Consumer sarama.ConsumerGroup
	Topics   string
}

func (f *FavoriteConsumerGroup) StartConsume(ctx context.Context, topics string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := f.Consumer.Consume(ctx, strings.Split(topics, ","), &handler); err != nil {
				logger.Log.Error("Error from consumer", zap.Any("error", err))
				continue
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
		}
	}()
	wg.Wait()
}

var handler FavoriteConsumerGroupHandler

type FavoriteConsumerGroupHandler struct {
}

func (FavoriteConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (FavoriteConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (FavoriteConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 获取消息
	for {
		select {
		case msg := <-claim.Messages():
			if err := processMessage(msg); err != nil {
				continue
			}
			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func processMessage(msg *sarama.ConsumerMessage) error {
	// 点赞kafka消息格式是key:videoId，userId:0?1，0表示取消点赞，1表示点赞
	// 1、校验、转换参数
	videoIdStr := string(msg.Key)
	videoId, err1 := strconv.Atoi(videoIdStr)
	c := context.Background()
	if strings.Contains(string(msg.Value), ":") {
		// value类型是userId:add或者userId:del
		valueList := strings.Split(string(msg.Value), ":")
		userIdStr := valueList[0]
		userId, err2 := strconv.Atoi(userIdStr)
		if len(valueList) != 2 || err1 != nil || err2 != nil {
			logger.Log.Error("消息格式错误")
			return errors.New("消息格式错误")
		}
		action := valueList[1]
		logger.Log.Info("receive msg", zap.Any("time", time.Now()), zap.Any("partition", msg.Partition), zap.Any("offset", msg.Partition),
			zap.Any("key", videoIdStr), zap.Any("userId", userId), zap.Any("videoId", videoId), zap.Any("action", action))
		// 2、操作redis
		if action == "del" {
			return delRedisFavoriteMess(c, userId, videoId)
		} else {
			return addRedisFavoriteMess(c, userId, videoId)
		}
	} else {
		// value类型是cnt，视频点赞数量
		cntStr := string(msg.Value)
		cnt, err2 := strconv.Atoi(cntStr)
		if err1 != nil || err2 != nil {
			logger.Log.Error("消息格式错误")
			return errors.New("消息格式错误")
		}
		logger.Log.Info("receive msg", zap.Any("time", time.Now()), zap.Any("partition", msg.Partition), zap.Any("offset", msg.Partition),
			zap.Any("key", videoIdStr), zap.Any("videoId", videoId), zap.Any("cnt", cnt))
		err := redis.SetVideoFavoriteCount(c, videoId, cnt)
		if err != nil {
			return errors.New("redis.SetVideoFavoriteCount 设置视频点赞数失败")
		}
	}
	return nil
}

func delRedisFavoriteMess(c context.Context, userId, videoId int) error {
	// 1) redis查看是否有点赞关系，如果没有不需要操作
	ok, err := redis.CheckIsFavorite(c, userId, videoId)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// 2) 如果有点赞关系
	// 先将redis关系删除，接着redis中视频点赞数-1
	err1 := redis.CancelFavorite(c, userId, videoId)
	err2 := redis.DelVideoFavoriteCount(c, videoId)
	if err1 != nil {
		return errors.New("redis.CancelFavorite 取消点赞操作失败")
	}
	if err2 != nil {
		_ = redis.AddFavorite(c, userId, videoId)
		return errors.New("redis.DelVideoFavoriteCount 减少视频点赞数失败")
	}
	return nil
}

func addRedisFavoriteMess(c context.Context, userId, videoId int) error {
	// 1) redis查看是否有点赞关系，如果有不需要操作
	ok, err := redis.CheckIsFavorite(c, userId, videoId)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	// 2) 如果没有点赞关系
	// 先将redis关系添加，接着redis中视频点赞数+1
	err1 := redis.AddFavorite(c, userId, videoId)
	err2 := redis.IncrVideoFavoriteCount(c, videoId)
	if err1 != nil {
		return errors.New("redis.AddFavorite 点赞操作失败")
	}
	if err2 != nil {
		_ = redis.CancelFavorite(c, userId, videoId)
		return errors.New("redis.IncrVideoFavoriteCount 增加视频点赞数失败")
	}
	return nil
}
