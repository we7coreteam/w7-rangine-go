package provider

import (
	"github.com/asaskevich/EventBus"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/src/core/console"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger"
)

type ProviderAbstract struct {
	ProviderInterface
	container     container.Container
	config        *viper.Viper
	console       *console.Console
	loggerFactory *logger.LoggerFactory
	event         EventBus.Bus
}

func (providerAbstract *ProviderAbstract) SetContainer(container container.Container) {
	providerAbstract.container = container
}

func (providerAbstract *ProviderAbstract) GetContainer() container.Container {
	return providerAbstract.container
}

func (providerAbstract *ProviderAbstract) SetConfig(config *viper.Viper) {
	providerAbstract.config = config
}

func (providerAbstract *ProviderAbstract) GetConfig() *viper.Viper {
	return providerAbstract.config
}

func (providerAbstract *ProviderAbstract) SetConsole(console *console.Console) {
	providerAbstract.console = console
}

func (providerAbstract *ProviderAbstract) GetConsole() *console.Console {
	return providerAbstract.console
}

func (providerAbstract *ProviderAbstract) SetLoggerFactory(loggerFactory *logger.LoggerFactory) {
	providerAbstract.loggerFactory = loggerFactory
}

func (providerAbstract *ProviderAbstract) GetLoggerFactory() *logger.LoggerFactory {
	return providerAbstract.loggerFactory
}

func (providerAbstract *ProviderAbstract) SetEvent(event EventBus.Bus) {
	providerAbstract.event = event
}

func (providerAbstract *ProviderAbstract) GetEvent() EventBus.Bus {
	return providerAbstract.event
}
