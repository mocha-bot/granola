package main

import (
	"os"
	"os/signal"
	"syscall"

	zLog "github.com/rs/zerolog/log"
)

func main() {
	cfg, err := GetConfig()
	if err != nil {
		zLog.Fatal().Err(err).Msg("Failed to get config")
	}

	redis, err := NewRedisInstance(&cfg.RedisConfig)
	if err != nil {
		zLog.Fatal().Err(err).Msg("Failed to create redis client")
	}

	zLog.Info().Msgf("Config: %+v", cfg)
	zLog.Info().Msgf("Redis client created %+v", redis)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	zLog.Info().Msg("Starting application")

	<-quit

	zLog.Info().Msg("Shutting down application")
}