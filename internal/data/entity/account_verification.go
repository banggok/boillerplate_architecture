package entity

import (
	"time"

	"github.com/banggok/boillerplate_architecture/internal/config/app"
	"github.com/banggok/boillerplate_architecture/internal/config/validator"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/password"
)

type AccountVerification interface {
	Entity
	AccountID() uint
	Account() Account
	VerificationType() VerificationType
	Token() string
	ExpiresAt() time.Time
	Verified() bool
}

type VerificationType string

const (
	EMAIL_VERIFICATION VerificationType = "EMAIL"
)

type accountVerificationImpl struct {
	entity
	accountID        uint
	account          Account
	verificationType VerificationType
	token            string
	expiresAt        time.Time
	verified         bool
}

// Account implements AccountVerification.
func (a *accountVerificationImpl) Account() Account {
	return a.account
}

// AccountID implements AccountVerification.
func (a *accountVerificationImpl) AccountID() uint {
	return a.accountID
}

// ExpiresAt implements AccountVerification.
func (a *accountVerificationImpl) ExpiresAt() time.Time {
	return a.expiresAt
}

// CreatedAt implements AccountVerification.
// Subtle: this method shadows the method (entity).CreatedAt of accountVerificationImpl.entity.
func (a *accountVerificationImpl) CreatedAt() time.Time {
	return a.createdAt
}

// Token implements AccountVerification.
func (a *accountVerificationImpl) Token() string {
	return a.token
}

// VerificationType implements AccountVerification.
func (a *accountVerificationImpl) VerificationType() VerificationType {
	return a.verificationType
}

// Verified implements AccountVerification.
func (a *accountVerificationImpl) Verified() bool {
	return a.verified
}

type newAccountVerificationParams struct {
	verificationType VerificationType `validate:"required"`
}

type accountVerificationData struct {
	verificationType VerificationType `validate:"required"`
	token            string           `validate:"required"`
	expiresAt        time.Time        `validate:"required"`
	verified         bool             `validate:"required"`
}

func NewAccountVerificationData(verificationType VerificationType, token string,
	expiresAt time.Time, verified bool) accountVerificationData {
	return accountVerificationData{
		verificationType: verificationType,
		token:            token,
		expiresAt:        expiresAt,
		verified:         verified,
	}
}

type makeAccountVerificationParams struct {
	metadata
	accountVerificationData
	accountVerificationAssoc
}

type accountVerificationAssoc struct {
	accountID uint `validate:"required"`
	account   Account
}

func NewAccountVerificationAssoc(accountID uint, account Account) accountVerificationAssoc {
	return accountVerificationAssoc{
		accountID: accountID,
		account:   account,
	}
}

func MakeAccountVerification(metadata metadata, verificationData accountVerificationData, assoc accountVerificationAssoc,
) (AccountVerification, error) {
	var res AccountVerification
	param := makeAccountVerificationParams{
		metadata:                 metadata,
		accountVerificationAssoc: assoc,
		accountVerificationData:  verificationData}

	if err := validator.Validate.Struct(param); err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.AccountVerificationUnprocessEntity,
			"failed to validate new account verification entity")
	}

	if assoc.account != nil && assoc.accountID != assoc.account.ID() {
		return nil, custom_errors.New(
			nil,
			custom_errors.AccountVerificationUnprocessEntity,
			"accountID and account mismatch")
	}

	res = &accountVerificationImpl{
		entity: entity{
			id:        param.id,
			createdAt: param.createdAt,
			updatedAt: param.updatedAt,
		},
		accountID:        param.accountID,
		account:          assoc.account,
		verificationType: param.verificationType,
		token:            param.token,
		expiresAt:        param.expiresAt,
		verified:         param.verified,
	}

	return res, nil
}

func NewAccountVerification(verificationType VerificationType,
	account Account) (AccountVerification, error) {
	var res *accountVerificationImpl
	param := newAccountVerificationParams{
		verificationType: verificationType,
	}

	if err := validator.Validate.Struct(param); err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.AccountVerificationUnprocessEntity,
			"failed to validate new account verification entity")
	}

	res = &accountVerificationImpl{
		verificationType: param.verificationType,
		account:          account,
	}

	token, err := password.GeneratePassword(16)
	if err != nil || token == nil {
		return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to generate password")
	}

	res.token = *token
	expiredDuration := app.AppConfig.ExpiredDuration.EmailVerification
	res.expiresAt = time.Now().UTC().Add(expiredDuration)

	if account != nil {
		res.accountID = account.ID()
	}
	return res, nil
}
