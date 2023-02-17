package app

import (
	"github.com/asaskevich/EventBus"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/src/components/database"
	"github.com/we7coreteam/w7-rangine-go/src/components/redis"
	"github.com/we7coreteam/w7-rangine-go/src/components/translator"
	"github.com/we7coreteam/w7-rangine-go/src/core/console"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger"
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
)

var GApp *App

type App struct {
	Name            string
	Version         string
	config          *viper.Viper
	container       container.Container
	loggerFactory   *logger.Factory
	event           EventBus.Bus
	providerManager *provider.Manager
	console         *console.Console
}

func NewApp() *App {
	app := &App{
		Name:    "rangine",
		Version: "1.0.0",
	}

	app.InitConfig()
	app.InitContainer()
	app.InitEvent()
	app.InitLoggerFactory()
	app.InitConsole()
	app.InitProviderManager()
	app.RegisterProviders()

	GApp = app

	return app
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
	app.loggerFactory = logger.NewLoggerFactory()

	var loggerConfigMap map[string]logger.Config
	err := app.config.Unmarshal(&loggerConfigMap)
	if err != nil {
		panic(err)
	}

	app.loggerFactory.Register(loggerConfigMap)
}

func (app *App) GetLoggerFactory() *logger.Factory {
	return app.loggerFactory
}

func (app *App) InitEvent() {
	app.event = EventBus.New()
}

func (app *App) GetEvent() EventBus.Bus {
	return app.event
}

func (app *App) InitProviderManager() {
	app.providerManager = provider.NewProviderManager(app.container, app.config, app.loggerFactory, app.event, app.console)
}

func (app *App) GetProviderManager() *provider.Manager {
	return app.providerManager
}

func (app *App) RegisterProviders() {
	app.providerManager.RegisterProvider(new(translator.Provider)).Register()
	app.providerManager.RegisterProvider(new(database.Provider)).Register()
	app.providerManager.RegisterProvider(new(redis.Provider)).Register()
}

func (app *App) InitConsole() {
	app.console = console.NewConsole()
}

func (app *App) GetConsole() *console.Console {
	return app.console
}

func (app *App) RunConsole() {
	app.console.Run()
}
