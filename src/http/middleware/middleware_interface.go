package middleware

import "github.com/gin-gonic/gin"

type Interface interface {
	Process(ctx *gin.Context)
}
