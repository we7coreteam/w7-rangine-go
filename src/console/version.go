package console

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

type VersionCommand struct {
	Abstract
}

func (versionCommand VersionCommand) GetName() string {
	return "version"
}

func (versionCommand VersionCommand) GetDescription() string {
	return "version"
}

func (versionCommand VersionCommand) Configure(command *cobra.Command) {
	command.Flags().BoolP("version", "v", true, "version")
}

func (versionCommand VersionCommand) Handle(cmd *cobra.Command, args []string) {
	color.Infoln("version: 1.0.0")
}
