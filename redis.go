package main

import (
	"github.com/go-redis/redis"
)

func NewRedisInstance(conf RedisConfig) (*redis.Client, error) {
	redisConf := &redis.Options{
		Addr:     conf.GetAddress(),
		Password: conf.Password,
		DB:       conf.DB,
	}

	rClient := redis.NewClient(redisConf)

	if err := rClient.Ping().Err(); err != nil {
		return nil, err
	}

	return rClient, nil
}
