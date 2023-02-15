package middleware

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExceptionRecoverLogger struct {
	io.Writer
}

func (exceptionRecoverLogger *ExceptionRecoverLogger) Write(p []byte) (n int, err error) {
	//exceptionRecoverLogger.HandlerExceptions.GetExceptionHandler().Reporter(exceptionRecoverLogger.HandlerExceptions.Logger, fmt.Errorf("%v", err), string(exceptionRecoverLogger.HandlerExceptions.Stack(4)))
	return len(p), nil
}

type ExceptionMiddleware struct {
	MiddlewareAbstract
}

func (exceptionMiddleware ExceptionMiddleware) Process(ctx *gin.Context) {
	gin.CustomRecoveryWithWriter(&ExceptionRecoverLogger{}, func(ctx *gin.Context, err interface{}) {
		if gin.Mode() == gin.DebugMode {
			exceptionMiddleware.JsonResponseWithError(ctx, "", http.StatusInternalServerError)
		} else {
			exceptionMiddleware.JsonResponseWithError(ctx, "系统内部错误", http.StatusInternalServerError)
		}
	})
}
