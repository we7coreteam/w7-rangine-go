package redis

import "github.com/we7coreteam/w7-rangine-go/src/core/provider"

type RedisProvider struct {
	provider.ProviderAbstract
}

func (redisProvider *RedisProvider) Register() {
	var redisConfigMap map[string]Config
	err := redisProvider.GetConfig().Unmarshal(&redisConfigMap)
	if err != nil {
		panic(err)
	}

	redisFactory := NewRedisFactory()
	redisFactory.Register(redisConfigMap)
}
