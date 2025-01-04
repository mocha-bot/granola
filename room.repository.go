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
	FetchRatingByRoomSerial(ctx context.Context, serial string) (rate RatingSummary, err error)

	UpsertDocument(ctx context.Context, index string, room Room) (err error)
}

type Repository interface {
	RoomRepository
}

type repository struct {
	db     *gorm.DB
	search meilisearch.ServiceManager
}

func NewRepository(db *gorm.DB, search meilisearch.ServiceManager) Repository {
	return &repository{
		db:     db,
		search: search,
	}
}

func (r *repository) FetchRoomBySerial(ctx context.Context, serial string) (room Room, err error) {
	rawSQL := `
		WITH current_slug AS (
				SELECT reference, slug
				FROM slug_version
				WHERE is_current = 1 AND type = 'room' AND deleted_at IS NULL
		)
		SELECT 
				r.serial, 
				r.name, 
				cs.slug, 
				r.description, 
				r.created_by, 
				COUNT(rc.channel_serial) AS total_channel, 
				0.0 AS rate,
				r.created_at
		FROM 
				room r
				LEFT JOIN room_channel rc ON r.serial = rc.room_serial
				LEFT JOIN reference_tag rt ON r.serial = rt.reference
				LEFT JOIN current_slug cs ON r.serial = cs.reference
		WHERE 
				r.serial = ?
		GROUP BY 
				r.serial, r.name, cs.slug, r.description, r.created_by
		LIMIT 1
	`

	err = r.db.WithContext(ctx).Raw(rawSQL, serial).Scan(&room).Error
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

func (r *repository) FetchRatingByRoomSerial(ctx context.Context, serial string) (rate RatingSummary, err error) {
	err = r.db.WithContext(ctx).Table("rating_summary rs").
		Select("average_rating").
		Where("rs.reference = ?", serial).
		First(&rate).
		Error

	if err != nil && err == gorm.ErrRecordNotFound {
		err = nil
	}

	return
}

func (r *repository) UpsertDocument(ctx context.Context, index string, room Room) (err error) {
	roomJSON, err := json.Marshal([]Room{room})
	if err != nil {
		return
	}

	zLog.Debug().Msgf("Room JSON: %s", roomJSON)

	taskInfo, err := r.search.Index(index).AddDocumentsWithContext(ctx, roomJSON)
	if err != nil {
		return
	}

	zLog.Info().Msgf("TaskID: %#v", taskInfo)
	return
}
