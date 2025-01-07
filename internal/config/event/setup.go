package eventconfig

import (
	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	valueobject "github.com/banggok/boillerplate_architecture/internal/data/entity/value_object"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/event"
	"github.com/banggok/boillerplate_architecture/internal/services"
	"gorm.io/gorm"
)

const (
	EMAIL_VERIFICATION_SUCCESS event.EventTopic = "email_verification_status"
)

var EventBus *event.EventBus

func Setup(db *gorm.DB) {
	EventBus = event.NewEventBus(db)

	EventBus.Subscribe(EMAIL_VERIFICATION_SUCCESS, func(event event.Event) {
		accountVerification := event.Data.(*entity.AccountVerification)

		if (*accountVerification).VerificationType() != valueobject.EMAIL_VERIFICATION {
			event.Response <- custom_errors.New(nil, custom_errors.AccountVerificationUnprocessEntity, "verification type mismatch")
		}

		// Simulate database update
		err := services.ServiceConFig.AccountVerification().Verify(EventBus.Ctx, accountVerification)
		event.Response <- err // Send the response back to the producer
	})

}
