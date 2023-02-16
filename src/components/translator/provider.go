package translator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
)

type TranslatorProvider struct {
	provider.ProviderAbstract
}

func (translatorProvider *TranslatorProvider) Register() {
	translatorProvider.GetConfig().SetDefault("app.lang", "zh")

	uni := ut.New(zh.New())
	translator, _ := uni.GetTranslator(translatorProvider.GetConfig().GetString("app.lang"))
	_ = zh_translations.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), translator)

	err := translatorProvider.GetContainer().NamedSingleton("translator", func() ut.Translator {
		return translator
	})
	if err != nil {
		panic(err)
	}
}
