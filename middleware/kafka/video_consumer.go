package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"tiktok-demo/common"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/dao/redis"
	"tiktok-demo/logger"
	"time"
)

// VideoConsumerGroup 视频相关的消费者
type VideoConsumerGroup struct {
	Consumer sarama.ConsumerGroup
	Topics   string
}

func (f *VideoConsumerGroup) StartConsume(ctx context.Context, topics string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := f.Consumer.Consume(ctx, strings.Split(topics, ","), &videoHandler); err != nil {
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

var videoHandler VideoConsumerGroupHandler

type VideoConsumerGroupHandler struct {
}

func (VideoConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (VideoConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (VideoConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 获取消息
	for {
		select {
		case msg := <-claim.Messages():
			if err := processVideoMessage(msg); err != nil {
				continue
			}
			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func processVideoMessage(msg *sarama.ConsumerMessage) error {
	// 发布视频kafka消息格式是key:userId，value: json格式的Video信息
	// 1、校验、转换参数
	userIdStr := string(msg.Key)
	userId, err1 := strconv.Atoi(userIdStr)
	c := context.Background()
	var video common.Video
	err2 := json.Unmarshal(msg.Value, &video)
	if err1 != nil || err2 != nil {
		logger.Log.Error("消息格式错误")
		return errors.New("消息格式错误")
	}

	logger.Log.Info("receive msg", zap.Any("time", time.Now()), zap.Any("partition", msg.Partition), zap.Any("offset", msg.Partition),
		zap.Any("key", userId), zap.Any("videoInfo", video))
	// 2、先操作mysql，将用户发布视频信息插入mysql
	videoId, err := mysql.InsertVideo(int(video.Author.Id), video.Title, video.PlayUrl, video.CoverUrl)
	if err != nil {
		logger.Log.Error("mysql.InsertVideo failed", zap.Any("error", err))
		return err
	}
	// 3、接着将对应关系加入redis，key: tiktok:publish:user:userid；value: videoid集合
	err = redis.AddVideo(c, userId, videoId)
	if err != nil {
		_ = mysql.DeleteVideo(videoId)
		logger.Log.Error("redis.AddVideo failed", zap.Any("error", err))
		return err
	}
	return nil
}
