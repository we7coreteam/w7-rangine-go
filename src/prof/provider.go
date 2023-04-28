package prof

import (
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go-support/src/server"
)

type Provider struct {
}

func (provider Provider) Register(config *viper.Viper, serverFactory server.Factory) {
	var serverConfig Config
	err := config.UnmarshalKey("server.prof", &serverConfig)
	if err != nil {
		panic(err)
	}

	serverFactory.RegisterServer(NewProfServer(serverConfig))
}
