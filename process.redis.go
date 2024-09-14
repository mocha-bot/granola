package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	zLog "github.com/rs/zerolog/log"
)

var (
	// store channel name and message in map to be used later
	PubSubLake = make(map[string]chan *redis.Message)
)

type PubSub interface {
	ListenRedisPubSub(ctx context.Context)
	ProcessMessage(ctx context.Context)
}

type pubSub struct {
	RedisClient *redis.Client

	Channel     string
	TotalWorker int

	UseCase UseCase
}

func NewPubSub(redisClient *redis.Client, channel string, workers int, useCase UseCase) PubSub {
	constructor := &pubSub{
		RedisClient: redisClient,
		Channel:     channel,
		TotalWorker: workers,
		UseCase:     useCase,
	}

	PubSubLake[channel] = make(chan *redis.Message, constructor.TotalWorker)

	return constructor
}

func (h *pubSub) ListenRedisPubSub(ctx context.Context) {
	pubsub := h.RedisClient.Subscribe(ctx, h.Channel)
	pubsub.Subscribe(ctx, h.Channel)

	defer func() {
		if err := pubsub.Unsubscribe(ctx, h.Channel); err != nil {
			zLog.Error().Err(err).Msg("Failed to unsubscribe")
		}
		if err := pubsub.Close(); err != nil {
			zLog.Error().Err(err).Msg("Failed to close")
		}
	}()

	zLog.Info().Msgf("Listening to channel: %s", h.Channel)

	for {
		select {
		case <-ctx.Done():
			zLog.Info().Msgf("Stopping listening to the channel: %s", h.Channel)
			return
		case msg, ok := <-pubsub.Channel():
			if !ok {
				return
			}
			zLog.Info().Msgf("Received message: %+v", msg)

			PubSubLake[h.Channel] <- msg
		}
	}
}

func (h *pubSub) ProcessMessage(ctx context.Context) {
	for idx := 1; idx <= h.TotalWorker; idx++ {
		zLog.Info().Msgf("Starting worker %d at channel: %s", idx, h.Channel)
		go func(i int) {
			for {
				select {
				case <-ctx.Done():
					zLog.Info().Msgf("Stopping worker %d", i)
					return
				case msg, ok := <-PubSubLake[h.Channel]:
					if !ok {
						zLog.Info().Msgf("Worker %d channel closed", i)
						return
					}
					zLog.Info().Msgf("Worker %d received message: %+v", i, msg)

					if msg.Channel != h.Channel {
						zLog.Error().Msgf("Worker %d received message from different channel: %s", i, msg.Channel)
						continue
					}

					// ! for now, only process the room_update
					room, err := h.UseCase.GetRoomBySerial(ctx, msg.Payload)
					if err != nil {
						zLog.Error().Err(err).Msgf("Failed to get room by serial: %s", msg.Payload)
						continue
					}

					zLog.Debug().Msgf("Worker %d processed room: %+v", i, room)

					err = h.UseCase.AddToDocument(ctx, "rooms", room)
					if err != nil {
						zLog.Error().Err(err).Msgf("Failed to add room to document: %+v", room)
						continue
					}

					zLog.Debug().Msgf("Worker %d added room to document: %+v", i, room)
				case <-time.After(100 * time.Millisecond):
					// ! to prevent blocking it makes the cpu usage high
				}
			}
		}(idx)
	}
}
