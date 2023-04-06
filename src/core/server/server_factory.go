package server

var servers = make(map[string]Interface)

func RegisterServer(server Interface) {
	servers[server.GetServerName()] = server
}

func GetServer(serverName string) Interface {
	server, exists := servers[serverName]
	if !exists {
		return nil
	}

	return server
}
