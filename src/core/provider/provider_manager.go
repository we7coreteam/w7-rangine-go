package provider

import (
	"github.com/golobby/container/v3/pkg/container"
)

type ProviderManager struct {
	Container container.Container
}

func (providerManager *ProviderManager) MakeProvider(abstract ProviderInterface) ProviderInterface {
	abstract.SetContainer(providerManager.Container)
	return abstract
}
