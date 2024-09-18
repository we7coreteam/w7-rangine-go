package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/we7coreteam/w7-rangine-go/v2/pkg/support/facade"
	"github.com/we7coreteam/w7-rangine-go/v2/src/http/controller/validator/bind"
	httperf "github.com/we7coreteam/w7-rangine-go/v2/src/http/error"
	"github.com/we7coreteam/w7-rangine-go/v2/src/http/response"
)

type Abstract struct {
	response.Response
}

func (abstract Abstract) TranslateValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); !ok {
		return errors.WithMessage(err, "参数数据格式错误")
	} else {
		errStr := ""
		for _, e := range validationErrors {
			errStr += e.Translate(facade.GetTranslator()) + ";"
		}

		return httperf.ValidateFail{
			ValidateErrs: validationErrors,
			Msg:          errStr,
		}
	}
}

func (abstract Abstract) Validate(ctx *gin.Context, requestData interface{}) bool {
	err := ctx.ShouldBindWith(requestData, bind.NewCompositeBind(ctx))
	if err != nil {
		abstract.JsonResponseWithServerError(ctx, abstract.TranslateValidationError(err))
		return false
	}

	return true
}
