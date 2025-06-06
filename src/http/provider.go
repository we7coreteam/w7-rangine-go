package http

import (
	"github.com/spf13/viper"
	support "github.com/we7coreteam/w7-rangine-go/v3/pkg/support/console"
	"github.com/we7coreteam/w7-rangine-go/v3/pkg/support/server"
	"github.com/we7coreteam/w7-rangine-go/v3/src/http/console"
	"github.com/we7coreteam/w7-rangine-go/v3/src/http/response"
	httpserver "github.com/we7coreteam/w7-rangine-go/v3/src/http/server"
)

type Provider struct {
	server *httpserver.Server
}

func (provider *Provider) Register(config *viper.Viper, consoleManager support.IConsole, serverManager server.IManager) *Provider {
	response.Env = config.GetString("app.env")

	httpServer := httpserver.NewHttpDefaultServer(config)
	provider.server = httpServer

	serverManager.RegisterServer(httpServer)

	consoleManager.RegisterCommand(&console.RouteListCommand{
		Server: httpServer,
	})

	return provider
}

func (provider *Provider) Export() *httpserver.Server {
	return provider.server
}
