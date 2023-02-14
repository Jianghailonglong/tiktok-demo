package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
	"tiktok-demo/conf"
	"tiktok-demo/logger"
)

var (
	client *redis.Client
)

func InitRedis() error {
	client = redis.NewClient(&redis.Options{
		Addr:     conf.Config.RedisConfig.Host + ":" + strconv.Itoa(conf.Config.RedisConfig.Port),
		Password: conf.Config.RedisConfig.Password,
		DB:       conf.Config.RedisConfig.DB,
		PoolSize: conf.Config.RedisConfig.PoolSize,
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Log.Error("connect redis failed", zap.Any("error", err.Error()))
		return err
	}
	logger.Log.Info("init redis success")
	return nil
}
