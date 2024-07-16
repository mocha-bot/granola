package main

import (
	"context"

	"github.com/redis/go-redis/v9"
	zLog "github.com/rs/zerolog/log"
)

var (
	// store channel name and message in map to be used later
	// TODO: process the messages
	PubSubLake = make(map[string][]*redis.Message)
)

func ListenRedisPubSub(ctx context.Context, rClient *redis.Client, channel string) {
	pubsub := rClient.Subscribe(ctx, channel)
	defer func() {
		if err := pubsub.Unsubscribe(ctx, channel); err != nil {
			zLog.Error().Err(err).Msg("Failed to unsubscribe")
		}
		if err := pubsub.Close(); err != nil {
			zLog.Error().Err(err).Msg("Failed to close")
		}
	}()

	zLog.Info().Msgf("Listening to channel: %s", channel)

	messageCh := pubsub.Channel()
	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(doneCh)
				return
			case msg := <-messageCh:
				zLog.Info().Msgf("Received message: %+v", msg)
				PubSubLake[channel] = append(PubSubLake[channel], msg)
			}
		}
	}()

	<-doneCh
	zLog.Info().Msg("Listener shut down cleanly")
}
