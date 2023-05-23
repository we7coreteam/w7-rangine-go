package console

import (
	"github.com/spf13/cobra"
	"os"
	"strings"
)

type RootCommand struct {
	Abstract
}

func (rootCommand RootCommand) GetName() string {
	return ""
}

func (rootCommand RootCommand) GetDescription() string {
	return ""
}

func (rootCommand RootCommand) Configure(command *cobra.Command) {
	command.PersistentFlags().StringP("config-file", "f", "", "Set the configuration file path")
	command.PersistentFlags().StringArrayP("env-var", "e", make([]string, 0), "Set environment variables")
}

func (rootCommand RootCommand) Handle(command *cobra.Command, args []string) {
	// 将配置项写入配置中
	configFile, _ := command.Flags().GetString("config-file")
	if configFile != "" {
		os.Setenv("RANGINE_CONFIG_FILE", configFile)
	}
	env, err := command.Flags().GetStringArray("env-var")
	if err == nil {
		for _, val := range env {
			if strings.Index(val, "=") >= 0 {
				varArr := strings.Split(val, "=")
				os.Setenv(varArr[0], varArr[1])
			}
		}
	}
}
