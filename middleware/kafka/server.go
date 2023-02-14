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
	FavoriteServerGroup.Consumer = newConsumerGroup(strings.Split(conf.Config.KafkaConfig.EndPoint, ","), FavoriteGroupId)
	if FavoriteServerGroup.Consumer == nil {
		return errors.New("create new consumer failed")
	}
	ctx := context.Background()
	go FavoriteServerGroup.StartConsume(ctx, FavoriteTopic)
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
