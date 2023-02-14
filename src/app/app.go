package app

import (
	"github.com/asaskevich/EventBus"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/src/core/exception"
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
)

type App struct {
	Name            string
	Config          *viper.Viper
	Container       container.Container
	ProviderManager *provider.ProviderManager
}

func NewApp(config *viper.Viper) *App {
	return &App{
		Name:   "rangine",
		Config: config,
	}
}

func (app *App) registerContainer() {
	app.Container = container.New()
}

func (app *App) initProviderManager() {
	err := app.Container.NamedSingleton("provider_manager", func() *provider.ProviderManager {
		return &provider.ProviderManager{
			Container: app.Container,
		}
	})
	if err != nil {
		panic(err)
	}

	var providerManager *provider.ProviderManager
	err = app.Container.NamedResolve(&providerManager, "provider_manager")
	if err != nil {
		panic(err)
	}

	app.ProviderManager = providerManager
}

func (app *App) Bootstrap() {
	app.registerContainer()
	app.initProviderManager()
}

func (app *App) registerEvent() {
	app.Event = EventBus.New()
}

func (app *App) registerValidation() {
	uni := ut.New(zh.New())
	lang := app.Config.App.Lang
	if lang == "" {
		lang = "zh"
	}

	app.Translator, _ = uni.GetTranslator(lang)
	_ = zh_translations.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), app.Translator)
}
func (app *App) registerExceptionHandler() {
	app.HandlerExceptions = &exception.HandlerExceptions{
		Logger: app.LoggerFactory.Channel("default"),
	}
}
