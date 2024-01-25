package prof

import (
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/pkg/support/server"
)

type Provider struct {
}

func (provider Provider) Register(config *viper.Viper, serverManager server.Manager) {
	var serverConfig Config
	err := config.UnmarshalKey("server.prof", &serverConfig)
	if err != nil {
		panic(err)
	}

	serverManager.RegisterServer(NewProfServer(serverConfig))
}
