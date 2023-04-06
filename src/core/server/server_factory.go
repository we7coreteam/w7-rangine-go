package server

var servers = make(map[string]Interface)

func RegisterServer(server Interface) {
	servers[server.GetServerName()] = server
}

func GetServers() map[string]Interface {
	return servers
}
