package response

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	errorhandler "github.com/we7coreteam/w7-rangine-go/src/core/err_handler"
	"net/http"
)

var Env = "release"

type ErrFormatter func(ctx *gin.Context, err error) error
type DataFormatter func(ctx *gin.Context, data any, err error, statusCode int) any

var responseErrFormatter ErrFormatter = func(ctx *gin.Context, err error) error {
	if errorhandler.Found(err) {
		errMsg := ""
		if errors.As(err, &errorhandler.ResponseError{}) {
			errMsg = err.Error()
		}
		if errMsg == "" {
			if Env == "debug" {
				errMsg = fmt.Sprintf("%+v", err)
			} else {
				errMsg = "系统内部错误"
			}
		}
		return errors.New(errMsg)
	}

	return nil
}

var responseDataFormatter DataFormatter = func(ctx *gin.Context, data any, err error, statusCode int) any {
	ret := map[string]interface{}{
		"data":  data,
		"code":  statusCode,
		"error": "",
	}
	if err != nil {
		ret["error"] = err.Error()
	}

	return ret
}

func SetResponseErrFormatter(formatter ErrFormatter) {
	responseErrFormatter = formatter
}

func SetResponseDataFormatter(formatter DataFormatter) {
	responseDataFormatter = formatter
}

func GetResponseErrFormatter() ErrFormatter {
	return responseErrFormatter
}

func GetResponseDataFormatter() DataFormatter {
	return responseDataFormatter
}

type Response struct {
}

func (response Response) JsonSuccessResponse(ctx *gin.Context) {
	response.JsonResponseWithoutError(ctx, "success")
}

func (response Response) JsonResponseWithoutError(ctx *gin.Context, data any) {
	response.JsonResponse(ctx, data, nil, http.StatusOK)
}

func (response Response) JsonResponseWithServerError(ctx *gin.Context, err error) {
	response.JsonResponseWithError(ctx, err, http.StatusInternalServerError)
}

func (response Response) JsonResponseWithError(ctx *gin.Context, err error, statusCode int) {
	response.JsonResponse(ctx, "", err, statusCode)
}

func (response Response) JsonResponse(ctx *gin.Context, data any, err error, statusCode int) {
	ctx.Set("response_data", data)
	ctx.Set("response_err", err)
	ctx.Set("response_code", statusCode)
	ctx.JSON(statusCode, GetResponseDataFormatter()(ctx, data, GetResponseErrFormatter()(ctx, err), statusCode))
	ctx.Abort()
}
