package console

import "github.com/spf13/cobra"

type Console struct {
	handler *cobra.Command
}

func NewConsole() *Console {
	return &Console{
		handler: &cobra.Command{Use: "rangine"},
	}
}

func (console *Console) RegisterCommand(command CommandInterface) {
	handler := &cobra.Command{
		Run: command.Handle,
	}
	command.Configure(handler)

	console.handler.AddCommand(handler)
}

func (console *Console) Run() {
	err := console.handler.Execute()
	if err != nil {
		panic(err)
	}
}
