package command

import (
	"github.com/spf13/cobra"
	"github.com/we7coreteam/w7-rangine-go/src/core/console"
	"github.com/we7coreteam/w7-rangine-go/src/http/server"
)

type ServerStartCommand struct {
	console.CommandAbstract
}

func (serverCommand *ServerStartCommand) GetName() string {
	return "server:start"
}

func (serverCommand *ServerStartCommand) GetDescription() string {
	return "server start"
}

func (serverCommand *ServerStartCommand) Handle(cmd *cobra.Command, args []string) {
	//app.GApp.GetConfig()
	server.GHttpServer.Start()
}
