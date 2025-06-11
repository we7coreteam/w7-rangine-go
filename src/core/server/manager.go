package server

import (
	"errors"
	"github.com/we7coreteam/w7-rangine-go/v3/pkg/support/server"
	"os"
	"strconv"
	"strings"
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

func (sm *Manager) getServersPidFilePath(servers []string) string {
	_, err := os.Stat("./runtime")
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir("./runtime", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return "./runtime/" + strings.Join(servers, "_") + ".pid"
}

func (sm *Manager) Start(servers []string) {
	for _, serverName := range servers {
		serverObj := sm.GetServer(serverName)
		if serverObj == nil {
			panic(errors.New("server " + serverName + " not found"))
		}

		go serverObj.Start()
	}

	pidPath := sm.getServersPidFilePath(servers)
	err := os.WriteFile(pidPath, []byte(strconv.Itoa(os.Getpid())), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (sm *Manager) Stop(servers []string) {
	pidPath := sm.getServersPidFilePath(servers)
	data, err := os.ReadFile(pidPath)
	if err != nil {
		panic(err)
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		panic(err)
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		panic(err)
	}
	err = process.Kill()
	if err != nil {
		panic(err)
	}

	_ = os.Remove(pidPath)
}
