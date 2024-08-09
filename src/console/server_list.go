package console

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/we7coreteam/w7-rangine-go/v2/pkg/support/facade"
	"os"
)

type ServerListCommand struct {
	Abstract
}

func NewServerListCommand() *ServerListCommand {
	return &ServerListCommand{}
}

func (serverCommand ServerListCommand) GetName() string {
	return "server:list"
}

func (serverCommand ServerListCommand) GetDescription() string {
	return "support server list"
}

func (serverCommand ServerListCommand) Handle(cmd *cobra.Command, args []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#"})
	t.AppendSeparator()
	for name, _ := range facade.GetServerManager().GetAllServer() {
		t.AppendRow([]interface{}{name})
	}
	t.Render()
}
