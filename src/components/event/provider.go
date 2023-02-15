package event

import (
	"github.com/asaskevich/EventBus"
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
	"github.com/we7coreteam/w7-rangine-go/src/global"
)

type EventProvider struct {
	provider.ProviderAbstract
}

func (eventProvider *EventProvider) Register() {
	global.G.Event = EventBus.New()
}
