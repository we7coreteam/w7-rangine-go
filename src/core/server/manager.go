package server

import (
	"errors"
	"github.com/we7coreteam/w7-rangine-go/v3/pkg/support/server"
)

type Manager struct {
	server.ManagerInterface
	servers map[string]server.ServerInterface
}

func NewDefaultServerManager() *Manager {
	return &Manager{
		servers: make(map[string]server.ServerInterface),
	}
}

func (sm *Manager) RegisterServer(server server.ServerInterface) {
	sm.servers[server.GetServerName()] = server
}

func (sm *Manager) GetAllServer() map[string]server.ServerInterface {
	return sm.servers
}

func (sm *Manager) GetServer(serverName string) server.ServerInterface {
	s, exists := sm.servers[serverName]
	if !exists {
		return nil
	}

	return s
}

func (sm *Manager) Start(servers []string) {
	for _, serverName := range servers {
		serverObj := sm.GetServer(serverName)
		if serverObj == nil {
			panic(errors.New("server " + serverName + " not found"))
		}

		go serverObj.Start()
	}
}
