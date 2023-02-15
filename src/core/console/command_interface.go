package console

import "github.com/spf13/cobra"

type CommandInterface interface {
	Configure(command *cobra.Command)
	Handle(cmd *cobra.Command, args []string)
}
