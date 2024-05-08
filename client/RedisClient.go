package client

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient interface {
	Save(key string, value string, ctx context.Context) error
	Get(key string, ctx context.Context) (string, error)
}

type redisClient struct {
	redis *redis.Client
}

type RedisClientCfg struct {
	RedisClient *redis.Client
}

func NewRedisClient(cfg RedisClientCfg) RedisClient {
	return &redisClient{
		redis: cfg.RedisClient,
	}
}

func (r redisClient) Save(key string, value string, ctx context.Context) error {
	// Use the Set method of the Redis client to save the key-value pair
	err := r.redis.Set(ctx, key, value, time.Hour*24).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r redisClient) Get(key string, ctx context.Context) (string, error) {
	// Use the Get method of the Redis client to retrieve the value for the specified key
	val, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
