package server

import "github.com/we7coreteam/w7-rangine-go-support/src/server"

type Manager struct {
	server.Manager
	servers map[string]server.Server
}

func NewDefaultServerManager() *Manager {
	return &Manager{
		servers: make(map[string]server.Server),
	}
}

func (sm *Manager) RegisterServer(server server.Server) {
	sm.servers[server.GetServerName()] = server
}

func (sm *Manager) GetAllServer() map[string]server.Server {
	return sm.servers
}

func (sm *Manager) GetServer(serverName string) server.Server {
	s, exists := sm.servers[serverName]
	if !exists {
		return nil
	}

	return s
}
