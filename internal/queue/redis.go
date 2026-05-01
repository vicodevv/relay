package queue

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(config RedisConfig) (*RedisClient, error) {
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       config.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	logrus.Info("Successfully connected to Redis")

	return &RedisClient{client}, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}
