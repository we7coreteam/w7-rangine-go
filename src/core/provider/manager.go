package provider

import (
	"github.com/asaskevich/EventBus"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/src/console"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger"
)

type Manager struct {
	container     container.Container
	config        *viper.Viper
	console       *console.Console
	loggerFactory *logger.Factory
	event         EventBus.Bus
}

func NewProviderManager(container container.Container, config *viper.Viper, loggerFactory *logger.Factory, event EventBus.Bus, console *console.Console) *Manager {
	return &Manager{
		container:     container,
		config:        config,
		console:       console,
		loggerFactory: loggerFactory,
		event:         event,
	}
}

func (manager *Manager) RegisterProvider(abstract Interface) Interface {
	abstract.SetContainer(manager.container)
	abstract.SetConfig(manager.config)
	abstract.SetConsole(manager.console)
	abstract.SetLoggerFactory(manager.loggerFactory)
	abstract.SetEvent(manager.event)
	return abstract
}
