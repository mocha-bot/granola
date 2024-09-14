package main

import (
	"context"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
)

func NewMeilisearchInstance(ctx context.Context, cfg MeilisearchConfig) (meilisearch.ServiceManager, error) {
	client := meilisearch.New(cfg.Host, meilisearch.WithAPIKey(cfg.MasterKey))

	if !client.IsHealthy() {
		return nil, fmt.Errorf("meilisearch is not healthy, please check the configuration")
	}

	return client, nil
}
