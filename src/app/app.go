package app

import (
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/src/components/database"
	"github.com/we7coreteam/w7-rangine-go/src/components/event"
	"github.com/we7coreteam/w7-rangine-go/src/components/logger"
	"github.com/we7coreteam/w7-rangine-go/src/components/redis"
	"github.com/we7coreteam/w7-rangine-go/src/components/validation"
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
)

type App struct {
	Name            string
	Config          *viper.Viper
	Container       container.Container
	ProviderManager *provider.ProviderManager
}

func NewApp() *App {
	app := &App{
		Name: "rangine",
	}

	app.InitConfig()
	app.InitContainer()
	app.InitProviderManager()
	app.RegisterProviders()

	return app
}

func (app *App) InitConfig() {
	conf := viper.New()
	conf.SetConfigFile("./.env")

	if err := conf.ReadInConfig(); err != nil {
		panic(err)
	}

	app.Config = conf
}

func (app *App) GetConfig() *viper.Viper {
	return app.Config
}

func (app *App) InitContainer() {
	app.Container = container.New()
}

func (app *App) InitProviderManager() {
	app.ProviderManager = &provider.ProviderManager{
		Container: app.Container,
		Config:    app.Config,
	}
}

func (app *App) RegisterProviders() {
	app.ProviderManager.RegisterProvider(new(logger.LoggerProvider)).Register()
	app.ProviderManager.RegisterProvider(new(event.EventProvider)).Register()
	app.ProviderManager.RegisterProvider(new(validation.ValidationProvider)).Register()
	app.ProviderManager.RegisterProvider(new(database.DatabaseProvider)).Register()
	app.ProviderManager.RegisterProvider(new(redis.RedisProvider)).Register()
}

func Run() {

}

//func (app *App) registerEvent() {
//	app.Event = EventBus.New()
//}
//
//func (app *App) registerValidation() {
//	uni := ut.New(zh.New())
//	lang := app.Config.App.Lang
//	if lang == "" {
//		lang = "zh"
//	}
//
//	app.Translator, _ = uni.GetTranslator(lang)
//	_ = zh_translations.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), app.Translator)
//}
