package http

import (
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
	"github.com/we7coreteam/w7-rangine-go/src/http/command"
	httpserver "github.com/we7coreteam/w7-rangine-go/src/http/server"
	"github.com/we7coreteam/w7-rangine-go/src/http/session"
)

type Provider struct {
	provider.Abstract
}

func (provider *Provider) Register() {
	provider.GetConsole().RegisterCommand(new(command.ServerStartCommand))

	err := provider.GetContainer().NamedSingleton("http-server", func() *httpserver.Server {
		var sessionConfig session.SessionConf
		var cookieConfig session.Cookie
		err := provider.GetConfig().UnmarshalKey("session", &sessionConfig)
		if err != nil {
			panic(err)
		}
		err = provider.GetConfig().UnmarshalKey("cookie", &cookieConfig)
		if err != nil {
			panic(err)
		}

		server := httpserver.NewHttpSerer(provider.GetConfig())
		server.Session = session.NewSession(sessionConfig, cookieConfig)
		return server
	})
	if err != nil {
		panic(err)
	}
}
