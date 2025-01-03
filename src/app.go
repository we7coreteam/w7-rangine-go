package app

import (
	"bytes"
	"github.com/asaskevich/EventBus"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/gookit/color"
	"github.com/spf13/viper"
	cons "github.com/we7coreteam/w7-rangine-go/v2/pkg/support/console"
	"github.com/we7coreteam/w7-rangine-go/v2/pkg/support/facade"
	log "github.com/we7coreteam/w7-rangine-go/v2/pkg/support/logger"
	"github.com/we7coreteam/w7-rangine-go/v2/pkg/support/server"
	"github.com/we7coreteam/w7-rangine-go/v2/src/components/database"
	"github.com/we7coreteam/w7-rangine-go/v2/src/components/redis"
	"github.com/we7coreteam/w7-rangine-go/v2/src/components/translator"
	"github.com/we7coreteam/w7-rangine-go/v2/src/console"
	"github.com/we7coreteam/w7-rangine-go/v2/src/core/helper"
	"github.com/we7coreteam/w7-rangine-go/v2/src/core/logger"
	sm "github.com/we7coreteam/w7-rangine-go/v2/src/core/server"
	"github.com/we7coreteam/w7-rangine-go/v2/src/prof"
	"go.uber.org/zap/exp/zapslog"
	"log/slog"
	"os"
)

var GApp *App

type Option struct {
	Name                string
	Version             string
	DefaultConfigLoader func(config *viper.Viper)
}

type App struct {
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
			panic(err)
		}

		app.config.SetConfigFile(customerConfigPath)
		content, err := helper.ParseConfigFileEnv(customerConfigPath)
		if err != nil {
			panic(err)
		}

		if err := app.config.MergeConfig(bytes.NewReader(content)); err != nil {
			panic(err)
		}
	}

	app.config.SetDefault("app.env", "release")
	app.config.SetDefault("app.lang", "zh")

	facade.Config = app.config
}

func (app *App) GetConfig() *viper.Viper {
	return app.config
}

func (app *App) InitContainer() {
	app.container = container.New()

	facade.Container = app.container
}

func (app *App) GetContainer() container.Container {
	return app.container
}

func (app *App) InitLoggerFactory() {
	factory := logger.NewLoggerFactory()

	var loggerConfigMap map[string]log.Config
	err := app.config.UnmarshalKey("log", &loggerConfigMap)
	if err != nil {
		panic(err)
	}

	factory.Register(loggerConfigMap)

	if _, exists := loggerConfigMap["default"]; exists {
		defaultLog, err := factory.Channel("default")
		if err != nil {
			panic(err)
		}

		defaultSlog := slog.New(zapslog.NewHandler(defaultLog.Core(), zapslog.WithName("default")))
		slog.SetDefault(defaultSlog)
	}

	app.loggerFactory = factory

	facade.LoggerFactory = app.loggerFactory
}

func (app *App) GetLoggerFactory() log.Factory {
	return app.loggerFactory
}

func (app *App) InitEvent() {
	app.event = EventBus.New()

	facade.Event = app.event
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

	facade.ServerManager = app.serverManager
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

	facade.Console = app.console
}

func (app *App) GetConsole() cons.Console {
	return app.console
}

func (app *App) RunConsole() {
	app.console.Run()
}
