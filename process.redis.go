package main

import (
	"context"

	"github.com/redis/go-redis/v9"
	zLog "github.com/rs/zerolog/log"
)

var (
	// store channel name and message in map to be used later
	PubSubLake = make(map[string]chan *redis.Message)
)

type RedisProcess interface {
	ListenRedisPubSub(ctx context.Context, channel string)
	ProcessMessage(ctx context.Context, channel string)
}

type redisProcess struct {
	RedisClient *redis.Client
	TotalWorker int

	UseCase UseCase
}

func NewRedisProcess(redisClient *redis.Client, workers int, useCase UseCase) RedisProcess {
	return &redisProcess{
		RedisClient: redisClient,
		TotalWorker: workers,
		UseCase:     useCase,
	}
}

func (h *redisProcess) ListenRedisPubSub(ctx context.Context, channel string) {
	pubsub := h.RedisClient.Subscribe(ctx, channel)
	pubsub.Subscribe(ctx, channel)

	defer func() {
		if err := pubsub.Unsubscribe(ctx, channel); err != nil {
			zLog.Error().Err(err).Msg("Failed to unsubscribe")
		}
		if err := pubsub.Close(); err != nil {
			zLog.Error().Err(err).Msg("Failed to close")
		}
	}()

	zLog.Info().Msgf("Listening to channel: %s", channel)

	PubSubLake[channel] = make(chan *redis.Message, h.TotalWorker)

	for {
		select {
		case <-ctx.Done():
			zLog.Info().Msgf("Stopping listening to the channel: %s", channel)
			return
		case msg, ok := <-pubsub.Channel():
			if !ok {
				return
			}
			zLog.Info().Msgf("Received message: %+v", msg)

			PubSubLake[channel] <- msg
		}
	}
}

func (h *redisProcess) ProcessMessage(ctx context.Context, channel string) {
	for idx := 1; idx <= h.TotalWorker; idx++ {
		zLog.Info().Msgf("Starting worker %d at channel: %s", idx, channel)
		go func(i int) {
			for {
				select {
				case <-ctx.Done():
					zLog.Info().Msgf("Stopping worker %d", i)
					return
				case msg, ok := <-PubSubLake[channel]:
					if !ok {
						zLog.Info().Msgf("Worker %d channel closed", i)
						return
					}
					zLog.Info().Msgf("Worker %d received message: %+v", i, msg)
				default:
					// Add default case to prevent blocking
				}
			}
		}(idx)
	}
}
