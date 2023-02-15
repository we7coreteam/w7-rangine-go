package provider

import (
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
)

type ProviderInterface interface {
	Register()
	SetContainer(container container.Container)
	SetConfig(config *viper.Viper)
}
