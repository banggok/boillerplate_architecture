package model

import "time"

type Metadata struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"` // Primary key, auto-incrementing
	CreatedAt time.Time `gorm:"autoCreateTime"`           // Automatically set when the record is created
	UpdatedAt time.Time `gorm:"autoUpdateTime"`           // Automatically set when the record is updated
}
