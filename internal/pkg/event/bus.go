package event

import (
	"sync"

	"github.com/banggok/boillerplate_architecture/internal/pkg/middleware/transaction"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EventTopic string

type Event struct {
	Name     EventTopic
	Data     interface{}
	Response chan error
}

type EventBus struct {
	sync.Mutex
	Ctx       *gin.Context
	listeners map[EventTopic][]func(Event)
}

func NewEventBus(db *gorm.DB) *EventBus {
	ctx := &gin.Context{}
	ctx.Set(transaction.DBTRANSACTION, db)
	return &EventBus{
		Ctx:       ctx,
		listeners: make(map[EventTopic][]func(Event)),
	}
}

func (bus *EventBus) Publish(event Event) {
	bus.Lock()
	defer bus.Unlock()

	if handlers, found := bus.listeners[event.Name]; found {
		for _, handler := range handlers {
			go handler(event)
		}
	}
}

func (bus *EventBus) Subscribe(eventName EventTopic, handler func(Event)) {
	bus.Lock()
	defer bus.Unlock()

	bus.listeners[eventName] = append(bus.listeners[eventName], handler)
}
