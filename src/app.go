package app

import (
	"github.com/asaskevich/EventBus"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/gookit/color"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go-support/src"
	cons "github.com/we7coreteam/w7-rangine-go-support/src/console"
	"github.com/we7coreteam/w7-rangine-go-support/src/facade"
	log "github.com/we7coreteam/w7-rangine-go-support/src/logger"
	"github.com/we7coreteam/w7-rangine-go-support/src/server"
	"github.com/we7coreteam/w7-rangine-go/src/components/database"
	"github.com/we7coreteam/w7-rangine-go/src/components/redis"
	"github.com/we7coreteam/w7-rangine-go/src/components/translator"
	"github.com/we7coreteam/w7-rangine-go/src/console"
	"github.com/we7coreteam/w7-rangine-go/src/core/config/encoding"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger"
	sm "github.com/we7coreteam/w7-rangine-go/src/core/server"
	"github.com/we7coreteam/w7-rangine-go/src/prof"
	"os"
)

var GApp *App

type Option struct {
	Name                string
	Version             string
	DefaultConfigLoader func(config *viper.Viper)
}

type App struct {
	support.App

	Name          string
	Version       string
	config        *viper.Viper
	container     container.Container
	loggerFactory log.Factory
	serverManager server.Manager
	event         EventBus.Bus
	console       cons.Console
}

func NewApp(option Option) *App {
	GApp = &App{
		Name:    "rangine",
		Version: "1.0.10",
	}
	GApp.ApplyOption(option)

	facade.SetApp(GApp)

	GApp.InitConsole()
	GApp.InitConfig(option)
	GApp.InitContainer()
	GApp.InitLoggerFactory()
	GApp.InitEvent()
	GApp.InitServerManager()
	GApp.RegisterProviders()

	return GApp
}

func (app *App) ApplyOption(option Option) {
	if option.Name != "" {
		app.Name = option.Name
	}
	if option.Version != "" {
		app.Version = option.Version
	}
}

func (app *App) InitConfig(option Option) {
	app.config = viper.New()
	_ = app.config.GetDecoderRegistry().RegisterDecoder("yaml", encoding.Codec{})
	app.config.AutomaticEnv()

	customerConfigPath := os.Getenv("RANGINE_CONFIG_FILE")
	if customerConfigPath == "" && option.DefaultConfigLoader == nil {
		color.Warnln("Warning: The configuration file is missing. Confirm whether the configuration file is required and specify it")
	}

	if option.DefaultConfigLoader != nil {
		option.DefaultConfigLoader(app.config)
	}

	if customerConfigPath != "" {
		_, err := os.Stat(customerConfigPath)
		if err != nil && os.IsNotExist(err) {
			return
		}

		app.config.SetConfigFile(customerConfigPath)
		if err := app.config.MergeInConfig(); err != nil {
			panic(err)
		}
	}

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

	factory.Register(loggerConfigMap)

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

func (app *App) RegisterProviders() {
	translator.Provider{}.Register(app.container)
	database.Provider{}.Register(app.config, app.loggerFactory, app.container)
	redis.Provider{}.Register(app.config, app.container)
	prof.Provider{}.Register(app.config, app.GetServerManager())
}

func (app *App) InitServerManager() {
	app.serverManager = sm.NewDefaultServerManager()
}

func (app *App) GetServerManager() server.Manager {
	return app.serverManager
}

func (app *App) InitConsole() {
	app.console = console.NewConsole()

	app.console.RegisterCommand(new(console.MakeModuleCommand))
	app.console.RegisterCommand(new(console.MakeProjectCommand))
	app.console.RegisterCommand(new(console.MakeModelCommand))
	app.console.RegisterCommand(new(console.MakeCmdCommand))
	app.console.RegisterCommand(&console.ServerStartCommand{
		Name: app.Name,
	})
	app.console.RegisterCommand(new(console.ServerStopCommand))
	app.console.RegisterCommand(new(console.ServerListCommand))
	app.console.RegisterCommand(&console.VersionCommand{
		Version: app.Version,
	})
}

func (app *App) GetConsole() cons.Console {
	return app.console
}

func (app *App) RunConsole() {
	app.console.Run()
}
