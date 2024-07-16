package main

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type AppConfig struct {
	Name     string `env:"APP_NAME" envDefault:"mochi"`
	Host     string `env:"APP_HOST" envDefault:"localhost"`
	Port     int    `env:"APP_PORT" envDefault:"8080"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
	TimeZone string `env:"TIME_ZONE" envDefault:"Asia/Jakarta"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

type MeilisearchConfig struct {
	Host      string `env:"MEILISEARCH_HOST" envDefault:"http://localhost:7700"`
	MasterKey string `env:"MEILISEARCH_MASTER_KEY" envDefault:""`
}

func (r RedisConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type Config struct {
	AppConfig
	RedisConfig
	MeilisearchConfig
}

func GetConfig() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return &cfg, err
}
