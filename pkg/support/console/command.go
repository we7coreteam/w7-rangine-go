package console

import "github.com/spf13/cobra"

type Command interface {
	GetName() string
	GetDescription() string
	Configure(cmd *cobra.Command)
	Handle(cmd *cobra.Command, args []string)
}
