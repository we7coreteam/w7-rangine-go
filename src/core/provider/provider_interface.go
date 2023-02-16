package provider

import (
	"github.com/asaskevich/EventBus"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/src/core/console"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger"
)

type ProviderInterface interface {
	Register()
	SetContainer(container container.Container)
	GetContainer() container.Container
	SetConfig(config *viper.Viper)
	GetConfig() *viper.Viper
	SetConsole(console *console.Console)
	GetConsole() *console.Console
	SetLoggerFactory(loggerFactory *logger.LoggerFactory)
	GetLoggerFactory() *logger.LoggerFactory
	SetEvent(event EventBus.Bus)
	GetEvent() EventBus.Bus
}
