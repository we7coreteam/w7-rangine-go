package command

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/we7coreteam/w7-rangine-go/src/core/console"
	"github.com/we7coreteam/w7-rangine-go/src/http/server"
	"os"
)

type RouteListCommand struct {
	console.CommandAbstract
}

func (listCommand *RouteListCommand) GetName() string {
	return "route:list"
}

func (listCommand *RouteListCommand) GetDescription() string {
	return "route list"
}

func (listCommand *RouteListCommand) Handle(cmd *cobra.Command, args []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Path", "Method", "Handler"})
	t.AppendSeparator()
	for index, route := range server.GHttpServer.GinEngine.Routes() {
		t.AppendRow([]interface{}{index, route.Path, route.Method, route.Handler})
	}
	t.Render()
}
