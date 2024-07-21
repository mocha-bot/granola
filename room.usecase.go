package main

import "context"

type UseCase interface {
	GetRoomBySerial(ctx context.Context, serial string) (room Room, err error)
	AddToDocument(ctx context.Context, index string, room Room) (err error)
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

	room.Tags = tags

	return
}

func (u *useCase) AddToDocument(ctx context.Context, index string, room Room) (err error) {
	err = u.repo.AddToDocument(ctx, index, room)
	return
}
