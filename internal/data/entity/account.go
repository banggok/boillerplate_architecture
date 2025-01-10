package entity

import (
	"github.com/banggok/boillerplate_architecture/internal/config/validator"
	valueobject "github.com/banggok/boillerplate_architecture/internal/data/entity/value_object"

	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/password"
)

type Account interface {
	iMetadata
	TenantId() uint
	Name() string
	Password() string
	Email() string
	Phone() string
	Tenant() Tenant
	accountAssoc
	accountAction
}

type verifiedAction string

const (
	EMAIL           verifiedAction = "email_verification"
	CHANGE_PASSWORD verifiedAction = "change_password"
	VERIFIED        verifiedAction = ""
)

type accountAction interface {
	VerificationAction() (*verifiedAction, error)
}

type accountAssoc interface {
	AccountVerifications() *[]AccountVerification
}

type accountImpl struct {
	metadataImpl
	name                 string
	email                string
	phone                string
	password             string
	tenantId             uint
	tenant               Tenant
	accountVerifications *[]AccountVerification
}

// VerificationAction implements Account.
func (a *accountImpl) VerificationAction() (*verifiedAction, error) {
	if a.accountVerifications == nil || len(*a.accountVerifications) == 0 {
		return nil, custom_errors.New(nil, custom_errors.InternalServerError, "accountVerifications was empty when call VerificationAction")
	}
	ret := VERIFIED
	for _, accountVerification := range *a.accountVerifications {
		if accountVerification.VerificationType() == valueobject.EMAIL_VERIFICATION &&
			!accountVerification.Verified() {
			ret = EMAIL
			return &ret, nil
		}

		if accountVerification.VerificationType() == valueobject.CHANGE_PASSWORD &&
			!accountVerification.Verified() {
			ret = CHANGE_PASSWORD
			return &ret, nil
		}
	}

	return &ret, nil
}

// AccountVerifications implements Account.
func (a *accountImpl) AccountVerifications() *[]AccountVerification {
	return a.accountVerifications
}

type accountIdentity struct {
	Name     string `validate:"required,alpha_space" example:"John Doe"`
	Email    string `validate:"required,email" example:"admin@example.com"`
	Phone    string `validate:"omitempty,e164" example:"+1234567890"`
	Password *string
}

func NewAccountIdentity(name, email, phone string, password *string) accountIdentity {
	return accountIdentity{
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: password,
	}
}

type accountTenant struct {
	TenantID uint
	Tenant   Tenant
}

func NewAccountTenant(tenantId uint, tenant Tenant) accountTenant {
	return accountTenant{
		TenantID: tenantId,
		Tenant:   tenant,
	}
}

type newAccountParams struct {
	accountIdentity
	Tenant Tenant
}

type makeAccountParams struct {
	metadata
	accountIdentity
	accountTenant
	Password string `validate:"required"`
}

func MakeAccount(
	metadata metadata,
	accountIdentity accountIdentity,
	accountTenant accountTenant,
	accountVerifications *[]AccountVerification) (Account, error) {
	var password string
	if accountIdentity.Password != nil {
		password = *accountIdentity.Password
	}
	params := makeAccountParams{
		metadata:        metadata,
		accountIdentity: accountIdentity,
		accountTenant:   accountTenant,
		Password:        password,
	}
	if err := validator.Validate.Struct(params); err != nil {
		return nil, custom_errors.New(
			err,
			custom_errors.AccountUnprocessEntity,
			"failed to validate make account entity")
	}

	return &accountImpl{
		metadataImpl: metadataImpl{
			id:        params.ID,
			createdAt: params.CreatedAt,
			updatedAt: params.UpdatedAt,
		},
		name:                 params.Name,
		email:                params.Email,
		phone:                params.Phone,
		tenantId:             params.TenantID,
		tenant:               params.Tenant,
		password:             params.Password,
		accountVerifications: accountVerifications,
	}, nil
}

func NewAccount(accountIdentity accountIdentity, tenant Tenant, accountVerifications *[]AccountVerification) (Account, *string, error) {
	params := newAccountParams{
		accountIdentity: accountIdentity,
		Tenant:          tenant,
	}
	if err := validator.Validate.Struct(params); err != nil {
		return nil, nil, custom_errors.New(
			err,
			custom_errors.AccountUnprocessEntity,
			"failed to validate new account entity")
	}

	plain, err := password.GeneratePassword(8)
	if err != nil || plain == nil {
		return nil, nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to generate password")
	}

	hashed, err := password.HashPassword(*plain)
	if err != nil || hashed == nil {
		return nil, nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to hash password")
	}

	account := accountImpl{
		name:                 params.Name,
		email:                params.Email,
		phone:                params.Phone,
		password:             *hashed,
		accountVerifications: accountVerifications,
	}

	if params.Tenant != nil {
		account.tenantId = params.Tenant.ID()
		account.tenant = params.Tenant
	}

	return &account, plain, nil
}

// TenantId retrieves the id of the tenant account.
func (a *accountImpl) TenantId() uint {
	return a.tenantId
}

// Name retrieves the name of the account.
func (a *accountImpl) Name() string {
	return a.name
}

// Password retrieves the hash password of the account.
func (a *accountImpl) Password() string {
	return a.password
}

// Email retrieves the email of the account.
func (a *accountImpl) Email() string {
	return a.email
}

// Phone retrieves the phone of the account.
func (a *accountImpl) Phone() string {
	return a.phone
}

// Tenant retrieves the tenant of the account.
func (a *accountImpl) Tenant() Tenant {
	return a.tenant
}
