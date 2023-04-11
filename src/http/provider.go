package http

import (
	"github.com/gin-gonic/gin"
	errorhandler "github.com/we7coreteam/w7-rangine-go/src/core/error"
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
	"github.com/we7coreteam/w7-rangine-go/src/core/server"
	"github.com/we7coreteam/w7-rangine-go/src/http/console"
	httperf "github.com/we7coreteam/w7-rangine-go/src/http/error"
	"github.com/we7coreteam/w7-rangine-go/src/http/response"
	httpserver "github.com/we7coreteam/w7-rangine-go/src/http/server"
	"net/http"
)

type Provider struct {
	provider.Abstract
}

func (provider *Provider) Register() {
	response.Env = provider.GetConfig().GetString("app.env")
	responseObj := response.Response{}

	httpServer := httpserver.NewHttpDefaultServer(provider.GetConfig())
	httpServer.Engine.NoRoute(func(context *gin.Context) {
		responseObj.JsonResponseWithError(context, httperf.NotFoundErr{
			Err: errorhandler.ResponseError{
				Msg: "Route not found, " + context.Request.URL.Path,
			},
		}, http.StatusNotFound)
	})
	httpServer.Engine.NoMethod(func(context *gin.Context) {
		responseObj.JsonResponseWithError(context, httperf.NotAllowErr{
			Err: errorhandler.ResponseError{
				Msg: "Route not allow, " + context.Request.URL.Path,
			},
		}, http.StatusNotFound)
	})

	server.RegisterServer(httpServer)

	provider.GetConsole().RegisterCommand(new(console.RouteListCommand))
}
