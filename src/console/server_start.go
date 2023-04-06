package console

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/src/core/server"
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
	color.Println("********************************************************************")

	for serverName, server := range server.GetServers() {
		go server.Start()

		color.Print(serverName + " | ")
		for key, val := range server.GetOptions() {
			color.Print(key + ": " + val + ",")
		}
		color.Println()
	}

	color.Println("********************************************************************")

	select {}
}
