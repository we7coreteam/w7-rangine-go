package http

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	support "github.com/we7coreteam/w7-rangine-go-support/src/console"
	"github.com/we7coreteam/w7-rangine-go-support/src/server"
	errorhandler "github.com/we7coreteam/w7-rangine-go/src/core/err_handler"
	"github.com/we7coreteam/w7-rangine-go/src/http/console"
	httperf "github.com/we7coreteam/w7-rangine-go/src/http/error"
	"github.com/we7coreteam/w7-rangine-go/src/http/response"
	httpserver "github.com/we7coreteam/w7-rangine-go/src/http/server"
	"net/http"
)

type Provider struct {
	server *httpserver.Server
}

func (provider *Provider) Register(config *viper.Viper, consoleManager support.Console, serverFactory server.Factory) *Provider {
	response.Env = config.GetString("app.env")
	responseObj := response.Response{}

	httpServer := httpserver.NewHttpDefaultServer(config)
	httpServer.Engine.HandleMethodNotAllowed = true
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
		}, http.StatusMethodNotAllowed)
	})
	provider.server = httpServer

	serverFactory.RegisterServer(httpServer)

	consoleManager.RegisterCommand(&console.RouteListCommand{
		Server: httpServer,
	})

	return provider
}

func (provider *Provider) Export() *httpserver.Server {
	return provider.server
}
