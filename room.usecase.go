package main

import (
	"context"
	"time"
)

type UseCase interface {
	GetRoomBySerial(ctx context.Context, serial string) (room Room, err error)
	UpsertDocument(ctx context.Context, index string, room Room) (err error)
}

type useCase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &useCase{
		repo: repo,
	}
}

func (u *useCase) GetRoomBySerial(ctx context.Context, serial string) (room Room, err error) {
	room, err = u.repo.FetchRoomBySerial(ctx, serial)
	if err != nil {
		return
	}

	tags, err := u.repo.FetchTagsByRoomSerial(ctx, serial)
	if err != nil {
		return
	}

	rateSummary, err := u.repo.FetchRatingByRoomSerial(ctx, serial)
	if err != nil {
		return
	}

	room.Tags = tags
	room.Rate = rateSummary
	room.UpdatedAt = time.Now()

	return
}

func (u *useCase) UpsertDocument(ctx context.Context, index string, room Room) (err error) {
	err = u.repo.UpsertDocument(ctx, index, room)
	return
}
