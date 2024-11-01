package translator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translational "github.com/go-playground/validator/v10/translations/zh"
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
)

type Provider struct {
}

func (provider Provider) Register(config *viper.Viper, container container.Container) {
	uni := ut.New(zh.New())
	translator, _ := uni.GetTranslator(config.GetString("app.lang"))
	_ = translational.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), translator)

	err := container.NamedSingleton("translator", func() ut.Translator {
		return translator
	})
	if err != nil {
		panic(err)
	}
}
