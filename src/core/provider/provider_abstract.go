package provider

import (
	"github.com/golobby/container/v3/pkg/container"
)

type ProviderAbstract struct {
	ProviderInterface
	Container container.Container
	Config    interface{}
}

func (providerAbstract *ProviderAbstract) SetContainer(container container.Container) {
	providerAbstract.Container = container
}

func (providerAbstract *ProviderAbstract) RegisterAutoPanic(name string, resolver interface{}) {
	err := providerAbstract.Container.NamedSingleton(name, resolver)
	if err != nil {
		panic(err)
	}
}
