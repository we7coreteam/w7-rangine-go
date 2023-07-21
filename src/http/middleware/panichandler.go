package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/we7coreteam/w7-rangine-go/src/core/err_handler"
	"github.com/we7coreteam/w7-rangine-go/src/http/response"
	"net/http"
)

func GetPanicHandlerMiddleware() gin.HandlerFunc {
	responseObj := response.Response{}
	return gin.CustomRecoveryWithWriter(nil, func(ctx *gin.Context, err any) {
		var recoverErr error
		if _, ok := err.(error); !ok {
			recoverErr = errors.New(fmt.Sprintf("%v", err))
		} else {
			recoverErr = err.(error)
		}

		responseObj.JsonResponseWithError(ctx, err_handler.Throw("系统内部错误", recoverErr), http.StatusInternalServerError)
	})
}
