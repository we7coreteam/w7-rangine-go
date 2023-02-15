package provider

import (
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/src/core/console"
)

type ProviderManager struct {
	container container.Container
	config    *viper.Viper
	console   *console.Console
}

func NewProviderManager(container container.Container, config *viper.Viper, console *console.Console) *ProviderManager {
	return &ProviderManager{
		container: container,
		config:    config,
		console:   console,
	}
}

func (providerManager *ProviderManager) RegisterProvider(abstract ProviderInterface) ProviderInterface {
	abstract.SetContainer(providerManager.container)
	abstract.SetConfig(providerManager.config)
	abstract.SetConsole(providerManager.console)
	return abstract
}
