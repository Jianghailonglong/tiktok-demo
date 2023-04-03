package kafka

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"strings"
	"tiktok-demo/conf"
	"tiktok-demo/logger"
)

// InitConsumerGroups 初始化消费者组
func InitConsumerGroups() error {
	ctx := context.Background()
	// 视频消费组
	VideoServerGroup.Consumer = newConsumerGroup(strings.Split(conf.Config.KafkaConfig.EndPoint, ","), VideoGroupId)
	if VideoServerGroup.Consumer == nil {
		return errors.New("create new consumer failed")
	}
	go VideoServerGroup.StartConsume(ctx, VideoTopic)

	// 点赞消费组
	f := NewFavoriteConsumerGroup()
	if f == nil {
		return errors.New("create new consumer failed")
	}
	go f.StartConsume(ctx, FavoriteTopic)

	// 聊天消费组
	ChatServerGroup.Consumer = newConsumerGroup(strings.Split(conf.Config.KafkaConfig.EndPoint, ","), ChatGroupId)
	if ChatServerGroup.Consumer == nil {
		return errors.New("create new consumer failed")
	}
	go ChatServerGroup.StartConsume(ctx, ChatTopic)

	// 关注消费组
	RelationServerGroup.Consumer = newConsumerGroup(strings.Split(conf.Config.KafkaConfig.EndPoint, ","), RelationGroupId)
	if RelationServerGroup.Consumer == nil {
		return errors.New("create new consumer failed")
	}
	go RelationServerGroup.StartConsume(ctx, RelationTopic)

	logger.Log.Info("init consumers success")
	return nil
}

func newConsumerGroup(addrs []string, groupId string) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	// Version 必须大于等于  V0_10_2_0
	config.Version = sarama.V0_10_2_1
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Retry.Max = 3
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategySticky}
	// 连接kafka
	group, err := sarama.NewConsumerGroup(addrs, groupId, config)
	if err != nil {
		logger.Log.Error("sarama.NewConsumerGroup failed", zap.Any("error", err))
	}
	return group
}
