package middleware

import (
	"github.com/gin-gonic/gin"
	errorhandler "github.com/we7coreteam/w7-rangine-go/src/core/error"
	"github.com/we7coreteam/w7-rangine-go/src/http/response"
	"io"
	"net/http"
)

type ErrRecoverLogger struct {
	io.Writer
	Reporter func(err error)
}

func (errRecoverLogger *ErrRecoverLogger) Write(p []byte) (n int, err error) {
	if errRecoverLogger.Reporter != nil {
		errRecoverLogger.Reporter(err)
	}

	return len(p), nil
}

type ErrHandlerMiddleware struct {
	Abstract
	response.Response
	Reporter func(err error)
}

func NewErrHandlerMiddleware(Reporter func(err error)) *ErrHandlerMiddleware {
	return &ErrHandlerMiddleware{
		Reporter: Reporter,
	}
}

func (handlerMiddleware ErrHandlerMiddleware) Process(ctx *gin.Context) {
	gin.CustomRecoveryWithWriter(&ErrRecoverLogger{
		Reporter: handlerMiddleware.Reporter,
	}, func(ctx *gin.Context, err interface{}) {
		handlerMiddleware.Response.JsonResponseWithError(ctx, errorhandler.Throw(err.(error).Error(), err.(error)), http.StatusInternalServerError)
	})
}
