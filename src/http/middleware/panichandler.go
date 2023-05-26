package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/we7coreteam/w7-rangine-go/src/http/response"
	"net/http"
)

func GetPanicHandlerMiddleware() gin.HandlerFunc {
	responseObj := response.Response{}
	return gin.CustomRecoveryWithWriter(nil, func(ctx *gin.Context, err any) {
		var recoverErr error
		if _, ok := err.(error); !ok {
			recoverErr = errors.New(err.(string))
		} else {
			recoverErr = err.(error)
		}

		responseObj.JsonResponseWithError(ctx, recoverErr, http.StatusInternalServerError)
	})
}
