package main

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AppConfig struct {
	Name     string `env:"APP_NAME" envDefault:"mochi"`
	Host     string `env:"APP_HOST" envDefault:"localhost"`
	Port     int    `env:"APP_PORT" envDefault:"8080"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
	TimeZone string `env:"TIME_ZONE" envDefault:"Asia/Jakarta"`
}

type DatabaseDialect string

const (
	DialectMySQL    DatabaseDialect = "mysql"
	DialectPostgres DatabaseDialect = "postgres"
)

type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST"`
	Port     int    `env:"DATABASE_PORT" envDefault:"3306"`
	Username string `env:"DATABASE_USERNAME" envDefault:"root"`
	Password string `env:"DATABASE_PASSWORD"`
	Schema   string `env:"DATABASE_SCHEMA"`
	Debug    bool   `env:"DATABASE_DEBUG" envDefault:"false"`
	Dialect  string `env:"DATABASE_DIALECT"`
}

func (d DatabaseConfig) GetDialector() gorm.Dialector {
	switch d.Dialect {
	case string(DialectMySQL):
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			d.Username,
			d.Password,
			d.Host,
			d.Port,
			d.Schema,
		)
		return mysql.Open(dsn)
	case string(DialectPostgres):
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
			d.Host,
			d.Username,
			d.Password,
			d.Schema,
			d.Port,
		)
		return postgres.Open(dsn)
	default:
		return nil
	}
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

type GoroutineConfig struct {
	Channel string `env:"CHANNEL" envDefault:""`
	Workers int    `env:"WORKERS" envDefault:"10"`
}

type Config struct {
	AppConfig
	DatabaseConfig
	RedisConfig
	MeilisearchConfig
	GoroutineConfig
}

func GetConfig() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return &cfg, err
}
