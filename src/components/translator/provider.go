package translator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translational "github.com/go-playground/validator/v10/translations/zh"
	"github.com/we7coreteam/w7-rangine-go-support/src/facade"
	"github.com/we7coreteam/w7-rangine-go-support/src/provider"
)

type Provider struct {
	provider.Abstract
}

func (provider *Provider) Register() {
	uni := ut.New(zh.New())
	translator, _ := uni.GetTranslator(facade.GetConfig().GetString("app.lang"))
	_ = translational.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), translator)

	err := facade.GetContainer().NamedSingleton("translator", func() ut.Translator {
		return translator
	})
	if err != nil {
		panic(err)
	}
}
