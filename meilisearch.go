package main

import (
	"fmt"

	"github.com/meilisearch/meilisearch-go"
)

func NewMeilisearchInstance(cfg MeilisearchConfig) (*meilisearch.Client, error) {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   cfg.Host,
		APIKey: cfg.MasterKey,
	})

	if !client.IsHealthy() {
		return nil, fmt.Errorf("meilisearch is not healthy, please check the configuration")
	}

	return client, nil
}
