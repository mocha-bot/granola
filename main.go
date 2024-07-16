package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis"
	"github.com/meilisearch/meilisearch-go"
	zLog "github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg, err := GetConfig()
	if err != nil {
		zLog.Fatal().Err(err).Msg("Failed to get config")
	}

	var redis *redis.Client
	var meilisearch *meilisearch.Client

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() (err error) {
		redis, err = NewRedisInstance(cfg.RedisConfig)
		if err != nil {
			err = fmt.Errorf("Failed to create redis client: %w", err)
		}
		return err
	})

	eg.Go(func() (err error) {
		meilisearch, err = NewMeilisearchInstance(cfg.MeilisearchConfig)
		if err != nil {
			err = fmt.Errorf("Failed to create meilisearch client: %w", err)
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		zLog.Fatal().Err(err).Msg("Failed to create clients")
	}

	zLog.Info().Msgf("Config: %+v", cfg)
	zLog.Info().Msgf("Redis client created %+v", redis)
	zLog.Info().Msgf("Meilisearch client created %+v", meilisearch)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	zLog.Info().Msg("Starting application")

	<-quit

	zLog.Info().Msg("Shutting down application")
}
