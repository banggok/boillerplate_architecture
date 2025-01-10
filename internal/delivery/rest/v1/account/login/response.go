package login

import (
	"time"

	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
)

type Response struct {
	ID           uint      `json:"id"`
	Action       string    `json:"action"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func transform(account entity.Account, accessToken string, refreshToken string) (*Response, error) {
	action, err := account.VerificationAction()
	if err != nil || action == nil {
		return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to get verification action")
	}
	return &Response{
		ID:           account.ID(),
		Action:       string(*action),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    account.CreatedAt(),
		UpdatedAt:    account.UpdatedAt(),
	}, nil
}
