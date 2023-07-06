package redis

import (
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/we7coreteam/w7-rangine-go/src/core/helper"
	"strconv"
	"strings"
	"sync"
)

type Factory struct {
	redisResolverMap map[string]func() redis.Cmdable
	redisMap         map[string]redis.Cmdable
	lock             sync.RWMutex
}

func NewRedisFactory() *Factory {
	return &Factory{
		redisMap:         make(map[string]redis.Cmdable),
		redisResolverMap: make(map[string]func() redis.Cmdable),
	}
}

func (factory *Factory) Channel(channel string) (redis.Cmdable, error) {
	factory.lock.RLock()
	redis, exists := factory.redisMap[channel]
	factory.lock.RUnlock()
	if exists {
		return redis, nil
	}

	factory.lock.Lock()
	defer factory.lock.Unlock()

	redis, exists = factory.redisMap[channel]
	if !exists {
		redisResolver, exists := factory.redisResolverMap[channel]
		if !exists {
			return nil, errors.New("redis channel " + channel + " not exists")
		}

		redis = redisResolver()
		factory.redisMap[channel] = redis
	}

	return redis, nil
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

func (factory *Factory) RegisterRedis(channel string, redisResolver func() redis.Cmdable) {
	factory.redisResolverMap[channel] = redisResolver
}

func (factory *Factory) Register(maps map[string]Config) {
	for key, value := range maps {
		func(channel string, config Config) {
			factory.RegisterRedis(channel, func() redis.Cmdable {
				fields := helper.ValidateAndGetErrFields(config)
				if len(fields) > 0 {
					panic("redis config error, channel: " + channel + ", fields: " + strings.Join(fields, ","))
				}

				return factory.MakeRedis(config)
			})
		}(key, value)
	}
}
