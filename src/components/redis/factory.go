package redis

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type Factory struct {
	channelMap map[string]redis.Cmdable
}

func NewRedisFactory() *Factory {
	return &Factory{
		channelMap: make(map[string]redis.Cmdable),
	}
}

func (factory *Factory) Channel(channel string) (redis.Cmdable, error) {
	cmdAble, exists := factory.channelMap[channel]
	if !exists {
		return nil, errors.New("redis channel " + channel + " not exists")
	}

	return cmdAble, nil
}

func (factory *Factory) MakeRedis(config Config) redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Username: config.Username,
		Password: config.Password,
		DB:       config.Db,
		PoolSize: config.PoolSize,
	})
}

func (factory *Factory) RegisterRedis(channel string, client redis.Cmdable) {
	factory.channelMap[channel] = client
}

func (factory *Factory) Register(maps map[string]Config) {
	for key, value := range maps {
		factory.RegisterRedis(key, factory.MakeRedis(value))
	}
}
