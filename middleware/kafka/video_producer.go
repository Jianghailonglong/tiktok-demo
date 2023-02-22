package kafka

import (
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"tiktok-demo/logger"
	"time"
)

// VideoProducer 视频相关的生产者
type VideoProducer struct {
	Client sarama.SyncProducer
	Topic  string
}

func (p *VideoProducer) SendMessage(key, value string) error {
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = p.Topic
	msg.Key = sarama.StringEncoder(key)
	msg.Value = sarama.StringEncoder(value)

	// 发送消息
	pid, offset, err := p.Client.SendMessage(msg)
	if err != nil {
		logger.Log.Error("send msg failed", zap.Any("error", err))
		return err
	}
	logger.Log.Info("send msg", zap.Any("time", time.Now()), zap.Any("partition", pid), zap.Any("offset", offset),
		zap.Any("key", msg.Key), zap.Any("value", msg.Value))
	return nil
}

func (p *VideoProducer) Close() error {
	err := p.Client.Close()
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	return nil
}
