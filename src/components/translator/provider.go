package translator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
	"github.com/we7coreteam/w7-rangine-go/src/global"
)

type TranslatorProvider struct {
	provider.ProviderAbstract
}

func (translatorProvider *TranslatorProvider) Register() {
	uni := ut.New(zh.New())
	lang := "zh"

	translator, _ := uni.GetTranslator(lang)
	_ = zh_translations.RegisterDefaultTranslations(binding.Validator.Engine().(*validator.Validate), translator)

	global.G.Translator = translator
}
