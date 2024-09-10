package main

import (
	"context"
	"encoding/json"

	"github.com/meilisearch/meilisearch-go"
	zLog "github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type RoomRepository interface {
	FetchRoomBySerial(ctx context.Context, serial string) (room Room, err error)
	FetchTagsByRoomSerial(ctx context.Context, serial string) (tags Tags, err error)

	AddToDocument(ctx context.Context, index string, room Room) (err error)
}

type Repository interface {
	RoomRepository
}

type repository struct {
	db     *gorm.DB
	search *meilisearch.Client
}

func NewRepository(db *gorm.DB, search *meilisearch.Client) Repository {
	return &repository{
		db:     db,
		search: search,
	}
}

func (r *repository) FetchRoomBySerial(ctx context.Context, serial string) (room Room, err error) {
	err = r.db.WithContext(ctx).Table("room r").
		Select("r.serial, r.name, r.description, r.created_by, COUNT(rc.channel_serial) AS total_channel, 0.0 AS rate, r.created_at, r.updated_at").
		Joins("LEFT JOIN room_channel rc ON r.serial = rc.room_serial").
		Joins("LEFT JOIN reference_tag rt ON r.serial = rt.reference").
		Where("r.serial = ?", serial).
		Group("r.serial, r.name, r.description, r.created_by").
		First(&room).Error

	return
}

func (r *repository) FetchTagsByRoomSerial(ctx context.Context, serial string) (tags Tags, err error) {
	err = r.db.WithContext(ctx).Table("reference_tag rt").
		Select("t.serial, t.name, t.description").
		Joins("JOIN tag t ON rt.tag_serial = t.serial").
		Where("rt.reference = ?", serial).
		Find(&tags).
		Error

	return
}

func (r *repository) AddToDocument(ctx context.Context, index string, room Room) (err error) {
	roomJSON, err := json.Marshal([]Room{room})
	if err != nil {
		return
	}

	zLog.Debug().Msgf("Room JSON: %s", roomJSON)

	taskInfo, err := r.search.Index(index).AddDocuments(roomJSON)
	if err != nil {
		return
	}

	zLog.Info().Msgf("TaskID: %#v", taskInfo)
	return
}
