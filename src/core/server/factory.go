package server

import "github.com/we7coreteam/w7-rangine-go-support/src/server"

type Factory struct {
	server.Factory
	servers map[string]server.Server
}

func NewDefaultServerFactory() *Factory {
	return &Factory{
		servers: make(map[string]server.Server),
	}
}

func (sf *Factory) RegisterServer(server server.Server) {
	sf.servers[server.GetServerName()] = server
}

func (sf *Factory) GetAllServer() map[string]server.Server {
	return sf.servers
}

func (sf *Factory) GetServer(serverName string) server.Server {
	s, exists := sf.servers[serverName]
	if !exists {
		return nil
	}

	return s
}
