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

func GetPanicHandlerMiddleware(reporter func(err error)) gin.HandlerFunc {
	responseObj := response.Response{}
	return gin.CustomRecoveryWithWriter(&ErrRecoverLogger{
		Reporter: reporter,
	}, func(ctx *gin.Context, err any) {
		var recoverErr error
		if _, ok := err.(error); !ok {
			recoverErr = errors.New(err.(string))
		} else {
			recoverErr = err.(error)
		}
		responseObj.JsonResponseWithError(ctx, errorhandler.Throw(recoverErr.Error(), recoverErr), http.StatusInternalServerError)
	})
}
