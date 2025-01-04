package main

import "time"

type Room struct {
	Serial       string         `json:"serial"`
	Name         string         `json:"name"`
	Slug         string         `json:"slug"`
	Description  string         `json:"description"`
	CreatedBy    string         `json:"created_by"`
	TotalChannel int            `json:"total_channel"`
	Tags         Tags           `json:"tags" gorm:"-"`
	Languages    Languages      `json:"languages" gorm:"-"`
	Rate         *RatingSummary `json:"rate" gorm:"-"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type Rooms []*Room

type Tag struct {
	Serial      string `json:"serial"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Tags []*Tag

type Languages []string

type RatingSummary struct {
	Reference string  `gorm:"column:reference" json:"reference"`
	Rate      float64 `gorm:"column:average_rating" json:"average_rating"`
	Count     int     `gorm:"column:rating_count" json:"rating_count"`
}

func (RatingSummary) TableName() string {
	return "rating_summary"
}
