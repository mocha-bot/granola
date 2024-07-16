package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabaseInstance(ctx context.Context, conf DatabaseConfig) (*gorm.DB, error) {
	gormConf := new(gorm.Config)
	gormConf.TranslateError = true

	if conf.Debug {
		gormConf.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Silent,
				Colorful:      true,
			},
		)
	}

	dialector := conf.GetDialector()
	if dialector == nil {
		return nil, fmt.Errorf("unsupported database dialect, %s", conf.Dialect)
	}

	instance, err := gorm.Open(dialector, gormConf)
	if err != nil {
		return nil, err
	}

	if conf.Debug {
		return instance.Debug(), nil
	}

	return instance, nil
}
