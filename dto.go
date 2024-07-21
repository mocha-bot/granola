package main

type Room struct {
	Serial       string    `json:"serial"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CreatedBy    string    `json:"created_by"`
	TotalChannel int       `json:"total_channel"`
	Tags         Tags      `json:"tags" gorm:"-"`
	Languages    Languages `json:"languages" gorm:"-"`
	Rate         float64   `json:"rate"`
}

type Rooms []*Room

type Tag struct {
	Serial      string `json:"serial"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Tags []*Tag

type Languages []string
