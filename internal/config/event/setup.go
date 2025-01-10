package eventconfig

import (
	"github.com/banggok/boillerplate_architecture/internal/data/entity"
	"github.com/banggok/boillerplate_architecture/internal/pkg/custom_errors"
	"github.com/banggok/boillerplate_architecture/internal/pkg/event"
	"github.com/banggok/boillerplate_architecture/internal/services"
	"gorm.io/gorm"
)

const (
	VERIFICATION_SUCCESS event.EventTopic = "verification_status"
)

var EventBus *event.EventBus

func Setup(db *gorm.DB) {
	EventBus = event.NewEventBus(db)

	EventBus.Subscribe(VERIFICATION_SUCCESS, func(event event.Event) {
		accountVerification, ok := event.Data.(*entity.AccountVerification)
		if !ok {
			event.Response <- custom_errors.New(nil, custom_errors.InternalServerError, "event data type invalid")
		}
		event.Response <- services.ServiceConFig.AccountVerification().Verify(EventBus.Ctx, accountVerification) // Send the response back to the producer
	})

}
