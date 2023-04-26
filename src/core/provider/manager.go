package provider

import (
	"github.com/we7coreteam/w7-rangine-go-support/src/provider"
)

type Manager struct {
}

func NewProviderManager() *Manager {
	return &Manager{}
}

func (manager *Manager) RegisterProvider(abstract provider.Provider) {
	abstract.Register()
}
