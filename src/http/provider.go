package http

import (
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
	"github.com/we7coreteam/w7-rangine-go/src/core/server"
	"github.com/we7coreteam/w7-rangine-go/src/http/console"
	httpserver "github.com/we7coreteam/w7-rangine-go/src/http/server"
)

type Provider struct {
	provider.Abstract
}

func (provider *Provider) Register() {
	server.RegisterServer(httpserver.NewHttpDefaultServer(provider.GetConfig()))

	provider.GetConsole().RegisterCommand(new(console.RouteListCommand))
}
