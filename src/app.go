package app

import (
	"github.com/asaskevich/EventBus"
	"github.com/golobby/container/v3/pkg/container"
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
	"github.com/we7coreteam/w7-rangine-go/src/core/logger"
	sm "github.com/we7coreteam/w7-rangine-go/src/core/server"
	"github.com/we7coreteam/w7-rangine-go/src/prof"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"os"
	"strings"
)

var GApp *App

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

func NewApp() *App {
	GApp = &App{
		Name:    "rangine",
		Version: "1.0.0",
	}

	facade.SetApp(GApp)

	GApp.InitConsole()
	GApp.InitConfig()
	GApp.InitContainer()
	GApp.InitLoggerFactory()
	GApp.InitEvent()
	GApp.InitServerManager()
	GApp.RegisterProviders()

	return GApp
}

func (app *App) InitConfig() {
	app.config = viper.New()
	app.config.AutomaticEnv()

	loadConfigFile := func(path string) {
		if path == "" {
			return
		}

		_, err := os.Stat(path)
		if err != nil && os.IsNotExist(err) {
			return
		}

		app.config.SetConfigFile(path)
		if err := app.config.MergeInConfig(); err != nil {
			panic(err)
		}
	}
	loadConfigFile("./config.yaml")
	loadConfigFile(os.Getenv("RANGINE_CONFIG_FILE"))

	var root = make(map[string]interface{})
	var buildMap func(path []string, value interface{}) map[string]interface{}
	buildMap = func(path []string, value interface{}) map[string]interface{} {
		if len(path) > 1 {
			return map[string]interface{}{
				path[0]: buildMap(path[1:], value),
			}
		}

		return map[string]interface{}{
			path[0]: value,
		}
	}
	for _, env := range os.Environ() {
		key, found := strings.CutPrefix(strings.Split(env, "=")[0], "RANGINE_")
		if found {
			path := strings.Split(strings.ToLower(key), ".")
			maps.Copy(root, buildMap(path, strings.Split(env, "=")[1]))
		}
	}
	err := app.config.MergeConfigMap(root)
	if err != nil {
		panic(err)
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
	app.console.RegisterCommand(&console.ServerStartCommand{
		Name: app.Name,
	})
	app.console.RegisterCommand(new(console.ServerStopCommand))
	app.console.RegisterCommand(new(console.ServerListCommand))
	app.console.RegisterCommand(new(console.VersionCommand))
}

func (app *App) GetConsole() cons.Console {
	return app.console
}

func (app *App) RunConsole() {
	app.console.Run()
}
