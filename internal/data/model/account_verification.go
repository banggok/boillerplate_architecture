package model

import (
	"time"

	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
)

// AccountVerification represents the account_verifications table
// This structure can be used with GORM for ORM mapping

type AccountVerification struct {
	Metadata
	AccountID uint `gorm:"not null;index"` // Foreign key referencing accounts
	Account   *Account
	Type      string    `gorm:"size:50;not null"`         // Type of verification (e.g., email, phone, etc.)
	Token     string    `gorm:"size:255;unique;not null"` // Unique token for verification
	ExpiresAt time.Time `gorm:"not null"`                 // Expiration timestamp
	Verified  bool      `gorm:"default:false"`            // Verification status
}

// TableName specifies the custom table name for GORM
func (AccountVerification) TableName() string {
	return "account_verifications"
}

func (AccountVerification) NotFoundError() custom_errors.ErrorCode {
	return custom_errors.AccountVerificationNotFound
}

func (a *AccountVerification) ToEntity() (entity.AccountVerification, error) {
	var accountEntity entity.Account
	if a.Account != nil {
		var err error
		accountEntity, err = a.Account.ToEntity()
		if err != nil {
			return nil, custom_errors.New(err, custom_errors.InternalServerError,
				"failed to convert account model to entity")
		}
	}
	return entity.MakeAccountVerification(
		entity.NewMetadata(a.ID, a.CreatedAt, a.UpdatedAt),
		entity.NewAccountVerificationData(entity.VerificationType(a.Type), a.Token, a.ExpiresAt, a.Verified),
		entity.NewAccountVerificationAssoc(a.AccountID, accountEntity))
}

func NewAccountVerification(accountVerificationEntity entity.AccountVerification) AccountVerification {
	var accountModel *Account
	if accountVerificationEntity.Account() != nil {
		toAccountModel := NewAccountModel(accountVerificationEntity.Account())
		accountModel = &toAccountModel
	}
	return AccountVerification{
		Metadata: Metadata{
			ID:        accountVerificationEntity.ID(),
			CreatedAt: accountVerificationEntity.CreatedAt().UTC(),
			UpdatedAt: accountVerificationEntity.UpdatedAt().UTC(),
		},
		AccountID: accountVerificationEntity.AccountID(),
		Account:   accountModel,
		Type:      string(accountVerificationEntity.VerificationType()),
		Token:     accountVerificationEntity.Token(),
		ExpiresAt: accountVerificationEntity.ExpiresAt().UTC(),
		Verified:  accountVerificationEntity.Verified(),
	}
}
