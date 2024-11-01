package http

import (
	"github.com/spf13/viper"
	appConsole "github.com/we7coreteam/w7-rangine-go/v2/src/console"
	"github.com/we7coreteam/w7-rangine-go/v2/src/core/server"
	"github.com/we7coreteam/w7-rangine-go/v2/src/http/console"
	"github.com/we7coreteam/w7-rangine-go/v2/src/http/response"
	httpserver "github.com/we7coreteam/w7-rangine-go/v2/src/http/server"
)

type Provider struct {
	server *httpserver.Server
}

func (provider *Provider) Register(config *viper.Viper, consoleManager *appConsole.Console, serverManager server.Manager) *Provider {
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
