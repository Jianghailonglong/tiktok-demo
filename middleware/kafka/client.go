package kafka

import (
	"errors"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"strings"
	"tiktok-demo/conf"
	"tiktok-demo/logger"
)

// InitProducers 初始化生产者
func InitProducers() error {
	addrStr := conf.Config.KafkaConfig.EndPoint
	addrs := strings.Split(addrStr, ",")

	VideoClient.Client = newProducer(addrs)
	if VideoClient.Client == nil {
		return errors.New("create new producer failed")
	}
	VideoClient.Topic = VideoTopic

	FavoriteClient.Client = newProducer(addrs)
	if FavoriteClient.Client == nil {
		return errors.New("create new producer failed")
	}
	FavoriteClient.Topic = FavoriteTopic

	ChatClient.Client = newProducer(addrs)
	if ChatClient.Client == nil {
		return errors.New("create new producer failed")
	}
	ChatClient.Topic = ChatTopic

	logger.Log.Info("init producers success")
	return nil
}

func newProducer(addrs []string) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Version = sarama.V0_10_2_1
	config.Producer.RequiredAcks = sarama.WaitForAll        // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewHashPartitioner // 对Key进行Hash，同样的Key每次都落到一个分区，这样消息是有序的
	config.Producer.Return.Successes = true                 // 成功交付的消息将在success channel返回
	config.Producer.Retry.Max = 3
	// 连接kafka
	client, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		logger.Log.Error("sarama.NewSyncProducer failed", zap.Any("error", err))
		return nil
	}

	return client
}
