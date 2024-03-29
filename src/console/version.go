package console

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

type VersionCommand struct {
	Abstract
	Version string
}

func (versionCommand VersionCommand) GetName() string {
	return "version"
}

func (versionCommand VersionCommand) GetDescription() string {
	return "version"
}

func (versionCommand VersionCommand) Configure(cmd *cobra.Command) {
	cmd.Flags().BoolP("version", "v", true, "version")
}

func (versionCommand VersionCommand) Handle(cmd *cobra.Command, args []string) {
	color.Infoln("version: " + versionCommand.Version)
}
