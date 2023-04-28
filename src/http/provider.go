package http

import (
	"github.com/spf13/viper"
	support "github.com/we7coreteam/w7-rangine-go-support/src/console"
	"github.com/we7coreteam/w7-rangine-go-support/src/server"
	"github.com/we7coreteam/w7-rangine-go/src/http/console"
	"github.com/we7coreteam/w7-rangine-go/src/http/response"
	httpserver "github.com/we7coreteam/w7-rangine-go/src/http/server"
)

type Provider struct {
	server *httpserver.Server
}

func (provider *Provider) Register(config *viper.Viper, consoleManager support.Console, serverFactory server.Factory) *Provider {
	response.Env = config.GetString("app.env")

	httpServer := httpserver.NewHttpDefaultServer(config)
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
