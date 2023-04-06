package provider

import (
	"github.com/asaskevich/EventBus"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/src/console"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger"
)

type Abstract struct {
	Interface
	container     container.Container
	config        *viper.Viper
	console       *console.Console
	loggerFactory *logger.Factory
	event         EventBus.Bus
	PackageName   string
}

func (abstract *Abstract) GetPackageName() string {
	return abstract.PackageName
}

func (abstract *Abstract) SetContainer(container container.Container) {
	abstract.container = container
}

func (abstract *Abstract) GetContainer() container.Container {
	return abstract.container
}

func (abstract *Abstract) SetConfig(config *viper.Viper) {
	abstract.config = config
}

func (abstract *Abstract) GetConfig() *viper.Viper {
	return abstract.config
}

func (abstract *Abstract) SetConsole(console *console.Console) {
	abstract.console = console
}

func (abstract *Abstract) GetConsole() *console.Console {
	return abstract.console
}

func (abstract *Abstract) SetLoggerFactory(loggerFactory *logger.Factory) {
	abstract.loggerFactory = loggerFactory
}

func (abstract *Abstract) GetLoggerFactory() *logger.Factory {
	return abstract.loggerFactory
}

func (abstract *Abstract) SetEvent(event EventBus.Bus) {
	abstract.event = event
}

func (abstract *Abstract) GetEvent() EventBus.Bus {
	return abstract.event
}
