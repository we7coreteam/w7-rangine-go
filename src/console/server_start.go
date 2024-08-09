package console

import (
	"errors"
	"github.com/gookit/color"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
	"github.com/we7coreteam/w7-rangine-go/v2/pkg/support/facade"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type ServerStartCommand struct {
	Abstract
	Name string
}

func (serverCommand ServerStartCommand) GetName() string {
	return "server:start"
}

func (serverCommand ServerStartCommand) Configure(cmd *cobra.Command) {
	cmd.Flags().BoolP("d", "d", false, "daemon start server")
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
		serverObj := facade.GetServerManager().GetServer(serverName)
		if serverObj == nil {
			panic(errors.New("server " + serverName + " not found"))
		}

		color.Print(serverName + " | ")
		for key, val := range serverObj.GetOptions() {
			color.Print(key + ": " + val + ",")
		}
		color.Println()
	}

	color.Println("********************************************************************")

	isDaemon, _ := cmd.Flags().GetBool("d")
	if isDaemon {
		ctx := &daemon.Context{
			WorkDir: "./",
			Umask:   027,
			Args:    []string{serverCommand.Name, "server:start"},
		}

		d, err := ctx.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		defer ctx.Release()
		if d != nil {
			return
		}

		return
	}

	facade.GetServerManager().Start(strings.Split(servers, "|"))

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	color.Println("Shutting down server...")
}
