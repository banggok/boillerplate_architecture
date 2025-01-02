package entity

import "time"

type Entity interface {
	GetID() uint
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}
type entity struct {
	id        uint      // id from database
	createdAt time.Time // timestamp from database
	updatedAt time.Time // timestamp from database
}

func (a *entity) GetID() uint {
	return a.id
}

func (a *entity) GetCreatedAt() time.Time {
	return a.createdAt
}

func (a *entity) GetUpdatedAt() time.Time {
	return a.updatedAt
}

type metadata struct {
	id        uint      `validate:"required"`
	createdAt time.Time `validate:"required"`
	updatedAt time.Time `validate:"required"`
}

func NewMetadata(id uint, createdAt, updatedAt time.Time) metadata {
	return metadata{
		id:        id,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
