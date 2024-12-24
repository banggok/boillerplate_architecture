package entity

import "time"

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
