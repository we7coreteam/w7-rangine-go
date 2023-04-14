package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	errorhandler "github.com/we7coreteam/w7-rangine-go/src/core/err_handler"
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

type PanicHandlerMiddleware struct {
	Abstract
	response.Response
	Reporter func(err error)
}

func NewPanicHandlerMiddleware(Reporter func(err error)) *PanicHandlerMiddleware {
	return &PanicHandlerMiddleware{
		Reporter: Reporter,
	}
}

func (handlerMiddleware PanicHandlerMiddleware) GetProcess() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(&ErrRecoverLogger{
		Reporter: handlerMiddleware.Reporter,
	}, func(ctx *gin.Context, err any) {
		var recoverErr error
		if _, ok := err.(error); !ok {
			recoverErr = errors.New(err.(string))
		} else {
			recoverErr = err.(error)
		}
		handlerMiddleware.Response.JsonResponseWithError(ctx, errorhandler.Throw(recoverErr.Error(), recoverErr), http.StatusInternalServerError)
	})
}
