package console

import "github.com/spf13/cobra"

type Interface interface {
	GetName() string
	GetDescription() string
	Configure(command *cobra.Command)
	Handle(cmd *cobra.Command, args []string)
}
