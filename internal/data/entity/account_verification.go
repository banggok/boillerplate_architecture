package entity

import (
	"time"

	"github.com/banggok/boillerplate_architecture/internal/config/app"
	"github.com/banggok/boillerplate_architecture/internal/config/validator"
	valueobject "github.com/banggok/boillerplate_architecture/internal/data/entity/value_object"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/password"
)

type AccountVerification interface {
	iMetadata
	AccountID() uint
	Account() Account
	VerificationType() valueobject.VerificationType
	Token() *string
	ExpiresAt() time.Time
	Verified() bool
	VerifiedSuccess()
}

var generateToken map[valueobject.VerificationType]bool = map[valueobject.VerificationType]bool{
	valueobject.EMAIL_VERIFICATION: true,
	valueobject.CHANGE_PASSWORD:    false,
}

var duration map[valueobject.VerificationType]time.Duration = map[valueobject.VerificationType]time.Duration{
	valueobject.EMAIL_VERIFICATION: app.AppConfig.ExpiredDuration.EmailVerification,
	valueobject.CHANGE_PASSWORD:    app.AppConfig.ExpiredDuration.ResetPasswordVerification,
}

type accountVerificationImpl struct {
	metadataImpl
	accountID        uint
	account          Account
	verificationType valueobject.VerificationType
	token            *string
	expiresAt        time.Time
	verified         bool
}

func (a *accountVerificationImpl) VerifiedSuccess() {
	a.verified = true
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
func (a *accountVerificationImpl) Token() *string {
	return a.token
}

// VerificationType implements AccountVerification.
func (a *accountVerificationImpl) VerificationType() valueobject.VerificationType {
	return a.verificationType
}

// Verified implements AccountVerification.
func (a *accountVerificationImpl) Verified() bool {
	return a.verified
}

type newAccountVerificationParams struct {
	VerificationType valueobject.VerificationType `validate:"required,verification_type"`
}

type accountVerificationData struct {
	VerificationType valueobject.VerificationType `validate:"required,verification_type"`
	Token            *string
	ExpiresAt        time.Time `validate:"required"`
	Verified         bool
}

func NewAccountVerificationData(verificationType valueobject.VerificationType, token *string,
	expiresAt time.Time, verified bool) accountVerificationData {
	return accountVerificationData{
		VerificationType: verificationType,
		Token:            token,
		ExpiresAt:        expiresAt,
		Verified:         verified,
	}
}

type makeAccountVerificationParams struct {
	metadata
	accountVerificationData
	accountVerificationAssoc
}

type accountVerificationAssoc struct {
	AccountID uint `validate:"required"`
	Account   Account
}

func NewAccountVerificationAssoc(accountID uint, account Account) accountVerificationAssoc {
	return accountVerificationAssoc{
		AccountID: accountID,
		Account:   account,
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

	if assoc.Account != nil && assoc.AccountID != assoc.Account.ID() {
		return nil, custom_errors.New(
			nil,
			custom_errors.AccountVerificationUnprocessEntity,
			"accountID and account mismatch")
	}

	res = &accountVerificationImpl{
		metadataImpl: metadataImpl{
			id:        param.ID,
			createdAt: param.CreatedAt,
			updatedAt: param.UpdatedAt,
		},
		accountID:        param.AccountID,
		account:          assoc.Account,
		verificationType: param.VerificationType,
		token:            param.Token,
		expiresAt:        param.ExpiresAt,
		verified:         param.Verified,
	}

	return res, nil
}

func NewAccountVerification(verificationType valueobject.VerificationType,
	account Account) (AccountVerification, error) {
	var res *accountVerificationImpl
	param := newAccountVerificationParams{
		VerificationType: verificationType,
	}

	if err := validator.Validate.Struct(param); err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.AccountVerificationUnprocessEntity,
			"failed to validate new account verification entity")
	}

	res = &accountVerificationImpl{
		verificationType: param.VerificationType,
		account:          account,
	}

	if generateToken[verificationType] {
		token, err := password.GeneratePassword(16)
		if err != nil || token == nil {
			return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to generate password")
		}

		res.token = token

	}

	res.expiresAt = time.Now().UTC().Add(duration[verificationType])

	if account != nil {
		res.accountID = account.ID()
	}
	return res, nil
}
