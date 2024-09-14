package main

import (
	_ "net/http/pprof"

	"context"
	"fmt"
	"net/http"
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

	if cfg.PPROF.IsEnabled {
		zLog.Info().Msg("Starting pprof...")
		go func() {
			zLog.Info().Msgf("pprof is now running on %s", cfg.PPROF.Address())
			err := http.ListenAndServe(cfg.PPROF.Address(), nil)
			if err != nil {
				zLog.Fatal().Err(err).Caller().Msg("error starting pprof")
			}
		}()
	}

	var redis *redis.Client
	var meilisearch meilisearch.ServiceManager
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

	repo := NewRepository(db, meilisearch)
	useCase := NewUseCase(repo)
	roomUpdate := NewPubSub(redis, cfg.GoroutineConfig.Channel, cfg.GoroutineConfig.Workers, useCase)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	zLog.Info().Msg("Starting application")

	var wg sync.WaitGroup
	delta := 1

	wg.Add(delta)
	go func() {
		defer wg.Done()
		roomUpdate.ListenRedisPubSub(ctx)
	}()

	wg.Add(delta)
	go func() {
		defer wg.Done()
		roomUpdate.ProcessMessage(ctx)
	}()

	<-ctx.Done()

	zLog.Info().Msg("Shutting down application")

	stop()
	wg.Wait()

	zLog.Info().Msg("Application shut down cleanly")
}
