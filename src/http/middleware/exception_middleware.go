package middleware

import (
	error_handler "github.com/we7coreteam/w7-rangine-go/src/core/error"
	"github.com/we7coreteam/w7-rangine-go/src/http/response"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExceptionRecoverLogger struct {
	io.Writer
	Reporter func(err error)
}

func (exceptionRecoverLogger *ExceptionRecoverLogger) Write(p []byte) (n int, err error) {
	if exceptionRecoverLogger.Reporter != nil {
		exceptionRecoverLogger.Reporter(err)
	}

	return len(p), nil
}

type ExceptionMiddleware struct {
	Abstract
	response.Response
	Reporter func(err error)
}

func (exceptionMiddleware ExceptionMiddleware) Process(ctx *gin.Context) {
	gin.CustomRecoveryWithWriter(&ExceptionRecoverLogger{
		Reporter: exceptionMiddleware.Reporter,
	}, func(ctx *gin.Context, err interface{}) {
		exceptionMiddleware.Response.JsonResponseWithError(ctx, error_handler.Throw(err.(error).Error(), err.(error)), http.StatusInternalServerError)
	})
}
