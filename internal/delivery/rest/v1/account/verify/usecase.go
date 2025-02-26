package verify

import (
	"time"

	eventconfig "github.com/banggok/boillerplate_architecture/internal/config/event"
	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	valueobject "github.com/banggok/boillerplate_architecture/internal/data/entity/value_object"
	"github.com/banggok/boillerplate_architecture/internal/delivery/rest/request"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/event"
	accountverification "github.com/banggok/boillerplate_architecture/internal/services/account_verification"
	"github.com/gin-gonic/gin"
)

type usecase interface {
	execute(ctx *gin.Context, request request.IRequest) (entity.Account, error)
}

type usecaseImpl struct {
	service accountverification.Service
}

var topicName = map[valueobject.VerificationType]event.EventTopic{
	valueobject.EMAIL_VERIFICATION: eventconfig.VERIFICATION_SUCCESS,
}

// Execute implements usecase.
func (u *usecaseImpl) execute(ctx *gin.Context, iRequest request.IRequest) (entity.Account, error) {
	request, ok := iRequest.(*Request)
	if !ok {
		return nil, custom_errors.New(nil, custom_errors.InternalServerError, "request invalid")
	}
	accountVerification, err := u.service.GetByTokenVerification(ctx, request.Token)
	if err != nil {
		return nil, custom_errors.New(err, custom_errors.InternalServerError, "token invalid")
	}

	if accountVerification.ExpiresAt().Before(time.Now()) {
		return nil, custom_errors.New(nil, custom_errors.AccountVerificationUnprocessEntity, "token expired")
	}

	if accountVerification == nil || accountVerification.Account() == nil {
		return nil, custom_errors.New(nil, custom_errors.InternalServerError, "account was empty")
	}

	if accountVerification.Verified() {
		return nil, custom_errors.New(nil, custom_errors.AccountVerificationConflictEntity, "token has been verified")
	}

	// Create a response channel to capture any errors from the subscribers
	responseChannel := make(chan error, 1)

	// Publish the event
	go func() {
		eventconfig.EventBus.Publish(event.Event{
			Name:     topicName[accountVerification.VerificationType()], // Topic of the event
			Data:     &accountVerification,                              // Data to send with the event
			Response: responseChannel,                                   // Channel to get the response
		})
	}()

	// Wait for the response from the subscriber
	err = <-responseChannel
	if err != nil {
		return nil, custom_errors.New(err, custom_errors.InternalServerError, "failed to update verified")
	}

	return accountVerification.Account(), nil
}

func newUsecase(service accountverification.Service) usecase {
	return &usecaseImpl{
		service: service,
	}
}
