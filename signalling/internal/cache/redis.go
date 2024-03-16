package cache

import (
	"context"

	"github.com/mvx-mnr-atomic/signalling/internal/logging"
	"github.com/mvx-mnr-atomic/signalling/internal/settings"
	"github.com/redis/go-redis/v9"
)

var logger = logging.GetLogger(nil)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(settings *settings.RedisSettings) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     settings.Address,
		Password: settings.Password,
		DB:       settings.DB,
	})
	return &RedisClient{
		client: client,
	}
}

func (rc *RedisClient) Get(key string) (string, error) {
	context := context.Background()
	val, err := rc.client.Get(context, key).Result()
	if err != nil {
		logger.Error("Error getting value from redis: ", err)
		return "", err
	}
	return val, nil
}

func (rc *RedisClient) Set(key string, value string) error {
	context := context.Background()
	err := rc.client.Set(context, key, value, 0).Err()
	if err != nil {
		logger.Error("Error setting value in redis: ", err)
		return err
	}
	return nil
}
