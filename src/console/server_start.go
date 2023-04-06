package console

import (
	"github.com/spf13/cobra"
	"github.com/we7coreteam/w7-rangine-go/src/core/server"
)

type ServerStartCommand struct {
	Abstract
}

func (serverCommand *ServerStartCommand) GetName() string {
	return "server:start"
}

func (serverCommand *ServerStartCommand) GetDescription() string {
	return "server start"
}

func (serverCommand *ServerStartCommand) Handle(cmd *cobra.Command, args []string) {
	for serverName, server := range server.GetServers() {
		go server.Start()
		cmd.OutOrStdout().Write(([]byte)("server " + serverName + " start"))
	}

	select {}
}
