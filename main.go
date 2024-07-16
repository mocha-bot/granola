package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/meilisearch/meilisearch-go"
	"github.com/redis/go-redis/v9"
	zLog "github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()

	cfg, err := GetConfig()
	if err != nil {
		zLog.Fatal().Err(err).Msg("Failed to get config")
	}

	zLog.Info().Msgf("Config: %+v", cfg)

	var redis *redis.Client
	var meilisearch *meilisearch.Client
	var db *gorm.DB

	eg, mCtx := errgroup.WithContext(ctx)

	eg.Go(func() (err error) {
		db, err = NewDatabaseInstance(mCtx, cfg.DatabaseConfig)
		if err != nil {
			err = fmt.Errorf("Failed to create database instance: %w", err)
		}
		return
	})

	eg.Go(func() (err error) {
		redis, err = NewRedisInstance(mCtx, cfg.RedisConfig)
		if err != nil {
			err = fmt.Errorf("Failed to create redis client: %w", err)
		}
		return
	})

	eg.Go(func() (err error) {
		meilisearch, err = NewMeilisearchInstance(mCtx, cfg.MeilisearchConfig)
		if err != nil {
			err = fmt.Errorf("Failed to create meilisearch client: %w", err)
		}
		return
	})

	if err := eg.Wait(); err != nil {
		zLog.Fatal().Err(err).Msg("Failed to create clients")
	}

	zLog.Info().Msgf("Database client created %+v", db)
	zLog.Info().Msgf("Redis client created %+v", redis)
	zLog.Info().Msgf("Meilisearch client created %+v", meilisearch)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	zLog.Info().Msg("Starting application")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		ListenRedisPubSub(ctx, redis, "channel")
	}()

	<-ctx.Done()

	zLog.Info().Msg("Shutting down application")

	stop()
	wg.Wait()

	zLog.Info().Msg("Application shut down cleanly")
}
