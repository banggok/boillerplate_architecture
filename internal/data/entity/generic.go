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
	ID        uint      `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
}

func NewMetadata(id uint, createdAt, updatedAt time.Time) metadata {
	return metadata{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
