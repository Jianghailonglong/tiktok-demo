package kafka

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/logger"
	"time"
)

// ChatConsumerGroup 聊天相关的消费者
type ChatConsumerGroup struct {
	Consumer sarama.ConsumerGroup
	Topics   string
}

func (f *ChatConsumerGroup) StartConsume(ctx context.Context, topics string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := f.Consumer.Consume(ctx, strings.Split(topics, ","), &chatHandler); err != nil {
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

var chatHandler ChatConsumerGroupHandler

type ChatConsumerGroupHandler struct {
}

func (ChatConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (ChatConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (ChatConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 获取消息
	for {
		select {
		case <-time.After(time.Second):
			// 每秒定时处理消息
			for msg := range claim.Messages() {
				if err := processChatMessage(msg); err != nil {
					continue
				}
				session.MarkMessage(msg, "")
			}
		case <-session.Context().Done():
			return nil
		}
	}
}

func processChatMessage(msg *sarama.ConsumerMessage) error {
	//  key:     toUserId，
	//	value:   userId:actionType:content
	// 1、校验、转换参数
	toUserIdStr := string(msg.Key)
	toUserId, err1 := strconv.Atoi(toUserIdStr)
	valueList := strings.Split(string(msg.Value), ":")
	userIdStr := valueList[0]
	userId, err2 := strconv.Atoi(userIdStr)
	if len(valueList) != 3 || err1 != nil || err2 != nil {
		logger.Log.Error("消息格式错误")
		return errors.New("消息格式错误")
	}
	actionType, _ := strconv.Atoi(valueList[1])
	content := valueList[2]
	logger.Log.Info("receive msg", zap.Any("time", time.Now()), zap.Any("partition", msg.Partition), zap.Any("offset", msg.Partition),
		zap.Any("key", toUserIdStr), zap.Any("userId", userId), zap.Any("actionType", actionType), zap.Any("content", content))

	// 2、操作mysql
	return mysql.InsertChatRecord(userId, toUserId, actionType, content)
}
