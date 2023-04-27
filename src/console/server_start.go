package console

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/we7coreteam/w7-rangine-go-support/src/facade"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type ServerStartCommand struct {
	Abstract
}

func (serverCommand ServerStartCommand) GetName() string {
	return "server:start"
}

func (serverCommand ServerStartCommand) GetDescription() string {
	return "server start"
}

func (serverCommand ServerStartCommand) Handle(cmd *cobra.Command, args []string) {
	servers := ""
	if len(args) == 0 {
		servers = facade.GetConfig().GetString("app.server")
	} else {
		servers = args[0]
	}
	if servers == "" {
		color.Errorln("please set the server to start")
		return
	}

	color.Println("********************************************************************")

	for _, serverName := range strings.Split(servers, "|") {
		serverObj := facade.GetServerFactory().GetServer(serverName)
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

	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	color.Println("Shutting down server...")
}
