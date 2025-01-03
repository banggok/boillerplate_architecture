package entity

import "time"

type Entity interface {
	ID() uint
	CreatedAt() time.Time
	UpdatedAt() time.Time
}
type entity struct {
	id        uint      // id from database
	createdAt time.Time // timestamp from database
	updatedAt time.Time // timestamp from database
}

func (a *entity) ID() uint {
	return a.id
}

func (a *entity) CreatedAt() time.Time {
	return a.createdAt
}

func (a *entity) UpdatedAt() time.Time {
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
