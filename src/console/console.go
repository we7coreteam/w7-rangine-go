package console

import (
	"github.com/spf13/cobra"
	command "github.com/we7coreteam/w7-rangine-go-support/src/console"
	"os"
)

type Console struct {
	handler *cobra.Command
}

func NewConsole() *Console {
	rootCommand := &RootCommand{}
	handler := &cobra.Command{
		Use:   rootCommand.GetName(),
		Short: rootCommand.GetDescription(),
	}
	rootCommand.Configure(handler)
	handler.FParseErrWhitelist = cobra.FParseErrWhitelist{
		UnknownFlags: true,
	}
	cmd, flags, _ := handler.Find(os.Args[1:])
	_ = cmd.ParseFlags(flags)
	rootCommand.Handle(cmd, flags)

	return &Console{
		handler: handler,
	}
}

func (console Console) RegisterCommand(cmd command.Command) {
	handler := &cobra.Command{
		Use:   cmd.GetName(),
		Short: cmd.GetDescription(),
		Run:   cmd.Handle,
	}
	cmd.Configure(handler)
	console.handler.AddCommand(handler)
}

func (console Console) GetHandler() *cobra.Command {
	return console.handler
}

func (console Console) Run() {
	_ = console.handler.Execute()
}
