package console

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/src/core/server"
	"strings"
)

type ServerStartCommand struct {
	Abstract

	config *viper.Viper
}

func NewServerStartCommand(config *viper.Viper) *ServerStartCommand {
	return &ServerStartCommand{
		config: config,
	}
}

func (serverCommand *ServerStartCommand) GetName() string {
	return "server:start"
}

func (serverCommand *ServerStartCommand) GetDescription() string {
	return "server start"
}

func (serverCommand *ServerStartCommand) Handle(cmd *cobra.Command, args []string) {
	servers := ""
	if len(args) == 0 {
		servers = serverCommand.config.GetString("app.server")
	} else {
		servers = args[0]
	}
	if servers == "" {
		color.Errorln("please set the server to start")
		return
	}

	color.Println("********************************************************************")

	for _, serverName := range strings.Split(servers, "|") {
		serverObj := server.GetServer(serverName)
		if serverObj == nil {
			color.Errorln("server " + serverName + " not exists!")
			return
		}

		go serverObj.Start()

		color.Print(serverName + " | ")
		for key, val := range serverObj.GetOptions() {
			color.Print(key + ": " + val + ",")
		}
		color.Println()
	}

	color.Println("********************************************************************")

	select {}
}