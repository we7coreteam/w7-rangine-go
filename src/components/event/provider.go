package event

import (
	"github.com/asaskevich/EventBus"
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
)

type EventProvider struct {
	provider.ProviderAbstract
}

func (eventProvider *EventProvider) Register() {
	EventBus.New()
}
