package console

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	command "github.com/we7coreteam/w7-rangine-go/v3/pkg/support/console"
	"github.com/we7coreteam/w7-rangine-go/v3/src/http/server"
	"os"
)

type RouteListCommand struct {
	command.CommandInterface
	Server *server.Server
}

func (listCommand RouteListCommand) Configure(cmd *cobra.Command) {

}

func (listCommand RouteListCommand) GetName() string {
	return "route:list"
}

func (listCommand RouteListCommand) GetDescription() string {
	return "route list"
}

func (listCommand RouteListCommand) Handle(cmd *cobra.Command, args []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Path", "Method", "Handler"})
	t.AppendSeparator()
	for index, route := range listCommand.Server.Engine.Routes() {
		t.AppendRow([]interface{}{index, route.Path, route.Method, route.Handler})
	}
	t.Render()
}
