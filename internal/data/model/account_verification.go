package model

import (
	"time"
)

// AccountVerification represents the account_verifications table
// This structure can be used with GORM for ORM mapping

type AccountVerification struct {
	ID        uint `gorm:"primaryKey;autoIncrement"` // Primary key
	AccountID uint `gorm:"not null;index"`           // Foreign key referencing accounts
	Account   *Account
	Type      string    `gorm:"size:50;not null"`         // Type of verification (e.g., email, phone, etc.)
	Token     string    `gorm:"size:255;unique;not null"` // Unique token for verification
	ExpiresAt time.Time `gorm:"not null"`                 // Expiration timestamp
	Verified  bool      `gorm:"default:false"`            // Verification status
	CreatedAt time.Time `gorm:"autoCreateTime"`           // Record creation timestamp
	UpdatedAt time.Time `gorm:"autoUpdateTime"`           // Last update timestamp
}

// TableName specifies the custom table name for GORM
func (AccountVerification) TableName() string {
	return "account_verifications"
}
