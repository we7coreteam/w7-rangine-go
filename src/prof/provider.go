package prof

import (
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/v3/pkg/support/server"
)

type Provider struct {
}

func (provider Provider) Register(config *viper.Viper, serverManager server.ManagerInterface) {
	var serverConfig Config
	err := config.UnmarshalKey("server.prof", &serverConfig)
	if err != nil {
		panic(err)
	}

	serverManager.RegisterServer(NewProfServer(serverConfig))
}
