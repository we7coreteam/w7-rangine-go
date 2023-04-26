package app

import (
	"github.com/asaskevich/EventBus"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go-support/src"
	cons "github.com/we7coreteam/w7-rangine-go-support/src/console"
	"github.com/we7coreteam/w7-rangine-go-support/src/facade"
	log "github.com/we7coreteam/w7-rangine-go-support/src/logger"
	"github.com/we7coreteam/w7-rangine-go/src/components/database"
	"github.com/we7coreteam/w7-rangine-go/src/components/redis"
	"github.com/we7coreteam/w7-rangine-go/src/components/translator"
	"github.com/we7coreteam/w7-rangine-go/src/console"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger"
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
	"github.com/we7coreteam/w7-rangine-go/src/prof"
	"go.uber.org/zap"
)

var GApp *App

type App struct {
	support.App
	Name            string
	Version         string
	config          *viper.Viper
	container       container.Container
	loggerFactory   log.Factory
	event           EventBus.Bus
	providerManager *provider.Manager
	console         cons.Console
}

func NewApp() *App {
	GApp = &App{
		Name:    "rangine",
		Version: "1.0.0",
	}

	facade.SetApp(GApp)

	GApp.InitConfig()
	GApp.InitContainer()
	GApp.InitLoggerFactory()
	GApp.InitEvent()
	GApp.InitConsole()
	GApp.InitProviderManager()

	return GApp
}

func (app *App) InitConfig() {
	conf := viper.New()
	conf.SetConfigFile("./config.yaml")

	if err := conf.ReadInConfig(); err != nil {
		panic(err)
	}

	app.config = conf

	app.config.SetDefault("app.env", "release")
	app.config.SetDefault("app.lang", "zh")
}

func (app *App) GetConfig() *viper.Viper {
	return app.config
}

func (app *App) InitContainer() {
	app.container = container.New()
}

func (app *App) GetContainer() container.Container {
	return app.container
}

func (app *App) InitLoggerFactory() {
	factory := logger.NewLoggerFactory()

	var loggerConfigMap map[string]logger.Config
	err := app.config.UnmarshalKey("log", &loggerConfigMap)
	if err != nil {
		panic(err)
	}

	for key, value := range loggerConfigMap {
		func(channel string, config logger.Config) {
			factory.RegisterLogger(channel, func() (*zap.Logger, error) {
				driver, err := factory.MakeDriver(config)
				if err != nil {
					return nil, err
				}
				return factory.MakeLogger(factory.ConvertLevel(config.Level), driver), nil
			})
		}(key, value)
	}

	app.loggerFactory = factory
}

func (app *App) GetLoggerFactory() log.Factory {
	return app.loggerFactory
}

func (app *App) InitEvent() {
	app.event = EventBus.New()
}

func (app *App) GetEvent() EventBus.Bus {
	return app.event
}

func (app *App) InitProviderManager() {
	app.providerManager = provider.NewProviderManager()

	app.providerManager.RegisterProvider(new(translator.Provider))
	app.providerManager.RegisterProvider(new(database.Provider))
	app.providerManager.RegisterProvider(new(redis.Provider))
	app.providerManager.RegisterProvider(new(prof.Provider))
}

func (app *App) GetProviderManager() *provider.Manager {
	return app.providerManager
}

func (app *App) InitConsole() {
	app.console = console.NewConsole()

	app.console.RegisterCommand(new(console.MakeModuleCommand))
	app.console.RegisterCommand(console.NewServerStartCommand(app.config))
	app.console.RegisterCommand(new(console.ServerListCommand))
	app.console.RegisterCommand(new(console.VersionCommand))
}

func (app *App) GetConsole() cons.Console {
	return app.console
}

func (app *App) RunConsole() {
	app.console.Run()
}
