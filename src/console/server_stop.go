package console

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/v2/src/core/server"
	"strings"
)

type ServerStopCommand struct {
	Abstract
	Config        *viper.Viper
	ServerManager *server.Manager
}

func (serverCommand ServerStopCommand) GetName() string {
	return "server:stop"
}

func (serverCommand ServerStopCommand) GetDescription() string {
	return "server stop"
}

func (serverCommand ServerStopCommand) Handle(cmd *cobra.Command, args []string) {
	servers := ""
	if len(args) == 0 {
		servers = serverCommand.Config.GetString("app.server")
	} else {
		servers = args[0]
	}
	if servers == "" {
		color.Errorln("please set the server to stop")
		return
	}

	serverCommand.ServerManager.Stop(strings.Split(servers, "|"))

	color.Successln("servers " + servers + " stop success")
}
