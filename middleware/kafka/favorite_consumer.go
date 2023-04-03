package kafka

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"tiktok-demo/conf"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/dao/redis"
	"tiktok-demo/logger"
	"time"
)

const (
	chanCount   = 10
	bufferCount = 1024
)

// FavoriteConsumerGroup 点赞相关的消费者
type FavoriteConsumerGroup struct {
	Consumer sarama.ConsumerGroup
}

func NewFavoriteConsumerGroup() *FavoriteConsumerGroup {
	f := &FavoriteConsumerGroup{
		Consumer: newConsumerGroup(strings.Split(conf.Config.KafkaConfig.EndPoint, ","), FavoriteGroupId),
	}

	if f.Consumer == nil {
		return nil
	}

	return f
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
	waiter   sync.WaitGroup
	msgsChan []chan *sarama.ConsumerMessage
}

func (f *FavoriteConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	handler.msgsChan = make([]chan *sarama.ConsumerMessage, chanCount)
	for i := 0; i < chanCount; i++ {
		ch := make(chan *sarama.ConsumerMessage, bufferCount)
		f.msgsChan[i] = ch
		f.waiter.Add(1)
		go f.consume(session, ch)
	}
	return nil
}

func (f *FavoriteConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	f.msgsChan = nil
	return nil
}

func (f *FavoriteConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 获取消息
	for {
		select {
		case msg := <-claim.Messages():
			key, err := strconv.Atoi(string(msg.Key))
			if err != nil {
				logger.Log.Error("msg.Key format is not correct")
				continue
			}
			f.msgsChan[key%chanCount] <- msg
		case <-session.Context().Done():
			return nil
		}
	}
}

func (f *FavoriteConsumerGroupHandler) consume(session sarama.ConsumerGroupSession, ch chan *sarama.ConsumerMessage) {
	defer f.waiter.Done()
	for {
		msg, _ := <-ch
		if err := processFavoriteMessage(msg); err != nil {
			continue
		}
		session.MarkMessage(msg, "")
	}
}

func processFavoriteMessage(msg *sarama.ConsumerMessage) error {
	// 点赞kafka消息格式是key:videoId，value: userId:del，删除，value: userId:add，添加
	// 1、校验、转换参数
	videoIdStr := string(msg.Key)
	videoId, err1 := strconv.Atoi(videoIdStr)
	c := context.Background()
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
		return procDelFavoriteMessage(c, userId, videoId)
	} else {
		return procAddFavoriteMessage(c, userId, videoId)
	}
}

func procDelFavoriteMessage(c context.Context, userId, videoId int) (err error) {
	// 1）先操作mysql
	// 先看是否原有点赞关系
	relation, _ := mysql.GetFavorite(userId, videoId)

	if relation != nil {
		err = mysql.UpdateFavorite(relation, mysql.UNFAVORITED)
	}
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	// 2）接着操作redis
	// 先将tiktok:favorite:user:userid对应的set集合删除videoid
	// 然后将tiktok:favorite:video:videoid对应的视频点赞数-1
	ok, err := redis.CheckIsFavorite(c, userId, videoId)
	if err != nil {
		return errors.New("redis.CheckIsFavorite 检查点赞关系失败")
	}
	if !ok {
		// 如果没有关系，后续不需要操作，直接返回
		return
	}
	err = redis.CancelFavorite(c, userId, videoId)
	if err != nil {
		return errors.New("redis.CancelFavorite 取消点赞操作失败")
	}
	err = redis.DecrVideoFavoriteCount(c, videoId)
	if err != nil {
		_ = redis.AddFavorite(c, userId, videoId) // 保持一致性
		_ = mysql.UpdateFavorite(relation, mysql.FAVORITED)
		return errors.New("redis.DecrVideoFavoriteCount 视频点赞数-1失败")
	}
	return
}

func procAddFavoriteMessage(c context.Context, userId, videoId int) (err error) {
	// 1）先操作mysql
	// 先看是否原有点赞关系
	relation, _ := mysql.GetFavorite(userId, videoId)

	if relation == nil {
		err = mysql.AddFavorite(userId, videoId)
	} else {
		err = mysql.UpdateFavorite(relation, mysql.FAVORITED)
	}
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	// 2）接着操作redis
	// 先将tiktok:favorite:user:userid对应的set集合添加videoid
	// 然后将tiktok:favorite:video:videoid对应的视频点赞数+1
	ok, err := redis.CheckIsFavorite(c, userId, videoId)
	if err != nil {
		return errors.New("redis.CheckIsFavorite 检查点赞关系失败")
	}
	if ok {
		// 如果有关系，后续不需要操作，直接返回
		return
	}
	err = redis.AddFavorite(c, userId, videoId)
	if err != nil {
		return errors.New("redis.AddFavorite 点赞操作失败")
	}
	err = redis.IncrVideoFavoriteCount(c, videoId)
	if err != nil {
		_ = redis.CancelFavorite(c, userId, videoId) // 保持一致性
		_ = mysql.UpdateFavorite(relation, mysql.UNFAVORITED)
		return errors.New("redis.IncrVideoFavoriteCount 视频点赞数+1失败")
	}
	return
}
