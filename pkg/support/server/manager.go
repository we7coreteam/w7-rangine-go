package server

type IManager interface {
	RegisterServer(server IServer)
	GetAllServer() map[string]IServer
	GetServer(serverName string) IServer
	Start(servers []string)
	Stop(servers []string)
}
