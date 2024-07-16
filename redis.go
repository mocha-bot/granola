package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisInstance(ctx context.Context, conf RedisConfig) (*redis.Client, error) {
	redisConf := &redis.Options{
		Addr:     conf.GetAddress(),
		Password: conf.Password,
		DB:       conf.DB,
	}

	rClient := redis.NewClient(redisConf)

	if err := rClient.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rClient, nil
}
