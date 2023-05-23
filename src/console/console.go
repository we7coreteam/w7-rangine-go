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
	rootCommandHandler := &RootCommand{}
	handler := &cobra.Command{
		Use:   rootCommandHandler.GetName(),
		Short: rootCommandHandler.GetDescription(),
	}
	rootCommandHandler.Configure(handler)
	cmd, flags, _ := handler.Find(os.Args[1:])
	_ = cmd.ParseFlags(flags)
	rootCommandHandler.Handle(cmd, flags)

	return &Console{
		handler: handler,
	}
}

func (console Console) RegisterCommand(command command.Command) {
	handler := &cobra.Command{
		Use:   command.GetName(),
		Short: command.GetDescription(),
		Run:   command.Handle,
	}
	command.Configure(handler)
	console.handler.AddCommand(handler)
}

func (console Console) Run() {
	_ = console.handler.Execute()
}
