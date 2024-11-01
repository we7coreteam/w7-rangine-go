package redis

import (
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
)

type Provider struct {
}

func (provider Provider) Register(config *viper.Viper, container container.Container) {
	var redisConfigMap map[string]Config
	err := config.UnmarshalKey("redis", &redisConfigMap)
	if err != nil {
		panic(err)
	}

	factory := NewRedisFactory()
	factory.Register(redisConfigMap)

	err = container.NamedSingleton("redis-factory", func() *Factory {
		return factory
	})
	if err != nil {
		panic(err)
	}
}
