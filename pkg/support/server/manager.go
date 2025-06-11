package server

type ManagerInterface interface {
	RegisterServer(server ServerInterface)
	GetAllServer() map[string]ServerInterface
	GetServer(serverName string) ServerInterface
	Start(servers []string)
	Stop(servers []string)
}
