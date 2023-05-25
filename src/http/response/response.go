package response

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	errorhandler "github.com/we7coreteam/w7-rangine-go/src/core/err_handler"
	"net/http"
)

var Env = "release"

type Formatter func(ctx *gin.Context, data any, err error, statusCode int) any

var responseFormatter Formatter = func(ctx *gin.Context, data any, err error, statusCode int) any {
	responseJson := map[string]interface{}{
		"data":  data,
		"code":  statusCode,
		"error": "",
	}

	if errorhandler.Found(err) {
		errMsg := ""
		if errors.As(err, &errorhandler.ResponseError{}) {
			errMsg = err.Error()
		}
		if errMsg == "" {
			if Env == "debug" {
				errMsg = err.Error()
				responseJson["err_strace"] = fmt.Sprintf("%+v \n ", err)
			} else {
				errMsg = "系统内部错误"
			}
		}
		responseJson["error"] = errMsg
	}

	return responseJson
}

func SetResponseFormatter(formatter Formatter) {
	responseFormatter = formatter
}

func GetResponseFormatter() Formatter {
	return responseFormatter
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
	ctx.JSON(statusCode, GetResponseFormatter()(ctx, data, err, statusCode))
	ctx.Abort()
}
