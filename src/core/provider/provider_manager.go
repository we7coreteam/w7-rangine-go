package provider

import (
	"github.com/asaskevich/EventBus"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/src/core/console"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger"
)

type ProviderManager struct {
	container     container.Container
	config        *viper.Viper
	console       *console.Console
	loggerFactory *logger.LoggerFactory
	event         EventBus.Bus
}

func NewProviderManager(container container.Container, config *viper.Viper, loggerFactory *logger.LoggerFactory, event EventBus.Bus, console *console.Console) *ProviderManager {
	return &ProviderManager{
		container:     container,
		config:        config,
		console:       console,
		loggerFactory: loggerFactory,
		event:         event,
	}
}

func (providerManager *ProviderManager) RegisterProvider(abstract ProviderInterface) ProviderInterface {
	abstract.SetContainer(providerManager.container)
	abstract.SetConfig(providerManager.config)
	abstract.SetConsole(providerManager.console)
	abstract.SetLoggerFactory(providerManager.loggerFactory)
	abstract.SetEvent(providerManager.event)
	return abstract
}
