package redis

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type Factory struct {
	channelMap map[string]*redis.Client
}

func NewRedisFactory() *Factory {
	return &Factory{
		channelMap: make(map[string]*redis.Client),
	}
}

func (factory *Factory) Channel(channel string) (*redis.Client, error) {
	redis, exists := factory.channelMap[channel]
	if !exists {
		return nil, errors.New("redis channel " + channel + " not exists")
	}

	return redis, nil
}

func (factory *Factory) MakeRedis(redisConfig Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host + ":" + strconv.Itoa(redisConfig.Port),
		Username: redisConfig.Username,
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
		PoolSize: redisConfig.PoolSize,
	})
}

func (factory *Factory) RegisterRedis(channel string, client *redis.Client) {
	factory.channelMap[channel] = client
}

func (factory *Factory) Register(maps map[string]Config) {
	for key, value := range maps {
		factory.RegisterRedis(key, factory.MakeRedis(value))
	}
}
