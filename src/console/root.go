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

func (rootCommand RootCommand) Configure(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("config-file", "f", "", "Set the configuration file path")
	cmd.PersistentFlags().StringArrayP("env-var", "e", make([]string, 0), "Set environment variables")
}

func (rootCommand RootCommand) Handle(cmd *cobra.Command, args []string) {
	// 将配置项写入配置中
	configFile, _ := cmd.Flags().GetString("config-file")
	if configFile != "" {
		_, err := os.Stat(configFile)
		if err != nil && os.IsNotExist(err) {
			panic(err)
		}

		os.Setenv("RANGINE_CONFIG_FILE", configFile)
	}
	env, err := cmd.Flags().GetStringArray("env-var")
	if err == nil {
		for _, val := range env {
			if strings.Index(val, "=") >= 0 {
				varArr := strings.Split(val, "=")
				os.Setenv(varArr[0], varArr[1])
			}
		}
	}
}
