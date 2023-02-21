package response

import (
	"fmt"
	app "github.com/we7coreteam/w7-rangine-go/src"
	error_handler "github.com/we7coreteam/w7-rangine-go/src/core/error"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Formatter func(ctx *gin.Context, data interface{}, error error, statusCode int) map[string]any

var responseFormatter Formatter = func(ctx *gin.Context, data interface{}, err error, statusCode int) map[string]any {
	errMsg := ""
	if error_handler.Found(err) {
		error_handler.Try(err).Catch(error_handler.ResponseError{}, func(err error) {
			errMsg = err.Error()
		}).Finally(func(err error) {
			if errMsg != "" {
				if gin.Mode() == gin.DebugMode {
					errMsg = fmt.Errorf("%+v", err).Error()
				} else {
					errMsg = "系统内部错误"
				}
			}
		})
	}

	return gin.H{
		"data": data,
		"code": statusCode,
		"msg":  errMsg,
	}
}

func SetResponseFormatter(formatter Formatter) {
	responseFormatter = formatter
}

func GetResponseFormatter() Formatter {
	return responseFormatter
}

type Response struct {
}

func (response *Response) JsonSuccessResponse(ctx *gin.Context) {
	response.JsonResponseWithoutError(ctx, "success")
}

func (response *Response) JsonResponseWithoutError(ctx *gin.Context, data interface{}) {
	response.JsonResponse(ctx, data, nil, http.StatusOK)
}

func (response *Response) JsonResponseWithServerError(ctx *gin.Context, err error) {
	response.JsonResponseWithError(ctx, err, http.StatusInternalServerError)
}

func (response *Response) JsonResponseWithError(ctx *gin.Context, err error, statusCode int) {
	response.JsonResponse(ctx, "", err, statusCode)
}

func (response *Response) JsonResponse(ctx *gin.Context, data interface{}, error error, statusCode int) {
	if error_handler.Found(error) {
		logger, _ := app.GApp.GetLoggerFactory().Channel("default")
		if logger != nil {
			logger.Debug(error.Error(), zap.Field{
				Type:      zapcore.ErrorType,
				Interface: error,
				Key:       "err",
			}, zap.Field{
				Type:      zapcore.UnknownType,
				Interface: ctx.Request,
				Key:       "request",
			})
		}
	}

	ctx.JSON(statusCode, GetResponseFormatter()(ctx, data, error, statusCode))
	ctx.Abort()
}
