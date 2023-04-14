package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/we7coreteam/w7-rangine-go/src/components/validator/bind"
	errorhandler "github.com/we7coreteam/w7-rangine-go/src/core/err_handler"
	"github.com/we7coreteam/w7-rangine-go/src/facade"
	httperf "github.com/we7coreteam/w7-rangine-go/src/http/error"
	"github.com/we7coreteam/w7-rangine-go/src/http/response"
)

type Abstract struct {
	response.Response
}

func (abstract Abstract) TranslateValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); !ok {
		return httperf.ValidateErr{
			Err: errorhandler.ResponseError{
				Msg: "参数数据格式错误",
			},
		}
	} else {
		errStr := ""
		for _, e := range validationErrors {
			errStr += e.Translate(facade.GetTranslator()) + ";"
		}

		return httperf.ValidateErr{
			ValidateErrs: validationErrors,
			Err: errorhandler.ResponseError{
				Msg: errStr,
			},
		}
	}
}

func (abstract Abstract) Validate(ctx *gin.Context, requestData interface{}) bool {
	err := ctx.ShouldBindWith(requestData, bind.NewCompositeBind(ctx.ContentType(), ctx.Params))
	if err != nil {
		abstract.JsonResponseWithServerError(ctx, abstract.TranslateValidationError(err))
		return false
	}

	return true
}
