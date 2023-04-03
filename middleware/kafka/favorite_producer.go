package kafka

import (
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"tiktok-demo/logger"
	"time"
)

// FavoriteProducer 点赞相关的生产者
type FavoriteProducer struct {
	Client sarama.SyncProducer
	Topic  string
}

func (p *FavoriteProducer) SendMessage(key, value string) error {
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

func (p *FavoriteProducer) SendMessages(vals map[string][]string) error {
	msgs := make([]*sarama.ProducerMessage, 0)
	for key, valSlice := range vals {
		for i := 0; i < len(valSlice); i++ {
			msg := &sarama.ProducerMessage{
				Topic: p.Topic,
				Key:   sarama.StringEncoder(key),
				Value: sarama.StringEncoder(valSlice[i]),
			}
			msgs = append(msgs, msg)
		}
	}

	// 发送消息
	err := p.Client.SendMessages(msgs)
	if err != nil {
		logger.Log.Error("send msgs failed", zap.Any("error", err))
		return err
	}
	logger.Log.Info("send msgs", zap.Any("time", time.Now()))
	return nil
}

func (p *FavoriteProducer) Close() error {
	err := p.Client.Close()
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	return nil
}
