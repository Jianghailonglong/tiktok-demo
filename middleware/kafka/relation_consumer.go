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
	"tiktok-demo/dao/redis"
	"tiktok-demo/logger"
	"time"
)

// RelationConsumerGroup 关注相关的消费者
type RelationConsumerGroup struct {
	Consumer sarama.ConsumerGroup
	Topics   string
}

func (f *RelationConsumerGroup) StartConsume(ctx context.Context, topics string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := f.Consumer.Consume(ctx, strings.Split(topics, ","), &relationHandler); err != nil {
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

var relationHandler RelationConsumerGroupHandler

type RelationConsumerGroupHandler struct {
}

func (RelationConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (RelationConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (RelationConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 获取消息
	for {
		select {
		case msg := <-claim.Messages():
			if err := processRelationMessage(msg); err != nil {
				continue
			}
			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func processRelationMessage(msg *sarama.ConsumerMessage) error {
	// 点赞kafka消息格式是key:userId，value: toUserId:del，删除，value: toUserId:add，添加
	// 1、校验、转换参数
	userIdStr := string(msg.Key)
	userId, err1 := strconv.Atoi(userIdStr)
	c := context.Background()
	// value类型是toUserId:add或者toUserId:del
	valueList := strings.Split(string(msg.Value), ":")
	toUserIdStr := valueList[0]
	toUserId, err2 := strconv.Atoi(toUserIdStr)
	if len(valueList) != 2 || err1 != nil || err2 != nil {
		logger.Log.Error("消息格式错误")
		return errors.New("消息格式错误")
	}
	action := valueList[1]
	logger.Log.Info("receive msg", zap.Any("time", time.Now()), zap.Any("partition", msg.Partition), zap.Any("offset", msg.Partition),
		zap.Any("key", userIdStr), zap.Any("userId", userId), zap.Any("toUserId", toUserId), zap.Any("action", action))
	// 2、操作redis
	if action == "del" {
		return procDelRelationMessage(c, userId, toUserId)
	} else {
		return procAddRelationMessage(c, userId, toUserId)
	}
}

func procDelRelationMessage(c context.Context, userId, toUserId int) (err error) {
	// 1）先操作mysql
	// 先看是否原有关注关系
	relation, _ := mysql.GetRelation(userId, toUserId)

	if relation != nil {
		err = mysql.UpdateRelation(relation, mysql.UNSUBSCRIBED)
	}

	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	// 2）接着操作redis
	// 先将tiktok:follow:user:userid对应的set集合删除to_userid
	// 然后将tiktok:follower:user:to_userid对应的set集合删除userid
	ok, err := redis.CheckIsRelation(c, userId, toUserId)
	if err != nil {
		return errors.New("redis.CheckIsRelation 检查关注关系失败")
	}
	if !ok {
		// 如果没有关系，后续不需要操作，直接返回
		return
	}

	err = redis.CancelRelation(c, userId, toUserId)
	if err != nil {
		_ = mysql.UpdateRelation(relation, mysql.SUBSCRIBED)
		return errors.New("redis.CancelRelation 取消关注操作失败")
	}
	return
}

func procAddRelationMessage(c context.Context, userId, toUserId int) (err error) {
	// 1）先操作mysql
	// 先看是否原有关注关系
	relation, _ := mysql.GetRelation(userId, toUserId)

	if relation == nil {
		err = mysql.AddRelation(userId, toUserId)
	} else {
		err = mysql.UpdateRelation(relation, mysql.SUBSCRIBED)
	}
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	// 2）接着操作redis
	// 一致性
	// 先将tiktok:follow:user:userid对应的set集合添加to_userid
	// 然后将tiktok:follower:user:to_userid对应的set集合添加userid
	ok, err := redis.CheckIsRelation(c, userId, toUserId)
	if err != nil {
		return errors.New("redis.CheckIsRelation 检查关注关系失败")
	}
	if ok {
		// 如果有关系，后续不需要操作，直接返回
		return
	}

	err = redis.AddRelation(c, userId, toUserId)
	if err != nil {
		_ = mysql.UpdateRelation(relation, mysql.UNSUBSCRIBED)
		return errors.New("redis.AddRelation 关注操作失败")
	}
	return
}
