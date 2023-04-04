package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	errorhandler "github.com/we7coreteam/w7-rangine-go/src/core/error"
	"github.com/we7coreteam/w7-rangine-go/src/facade"
	httperf "github.com/we7coreteam/w7-rangine-go/src/http/error"
	"github.com/we7coreteam/w7-rangine-go/src/http/response"
)

type Abstract struct {
	response.Response
}

func (abstract Abstract) TranslateValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); !ok {
		return err
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

func (abstract Abstract) ValidateFormPost(ctx *gin.Context, request interface{}) bool {
	err := ctx.ShouldBind(request)
	if err != nil {
		abstract.JsonResponseWithServerError(ctx, abstract.TranslateValidationError(err))
		return false
	}

	return true
}

func (abstract Abstract) ValidateQuery(ctx *gin.Context, request interface{}) bool {
	err := ctx.ShouldBindQuery(request)
	if err != nil {
		abstract.JsonResponseWithServerError(ctx, abstract.TranslateValidationError(err))
		return false
	}

	return true
}

func (abstract Abstract) ValidateUri(ctx *gin.Context, request interface{}) bool {
	err := ctx.ShouldBindUri(request)
	if err != nil {
		abstract.JsonResponseWithServerError(ctx, abstract.TranslateValidationError(err))
		return false
	}

	return true
}
