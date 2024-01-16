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
	redisResolverMap map[string]func() (redis.Cmdable, error)
	redisMap         map[string]redis.Cmdable
	lock             sync.RWMutex
}

func NewRedisFactory() *Factory {
	return &Factory{
		redisMap:         make(map[string]redis.Cmdable),
		redisResolverMap: make(map[string]func() (redis.Cmdable, error)),
	}
}

func (factory *Factory) Channel(channel string) (redis.Cmdable, error) {
	factory.lock.RLock()
	redisHandler, exists := factory.redisMap[channel]
	factory.lock.RUnlock()
	if exists {
		return redisHandler, nil
	}

	factory.lock.Lock()
	defer factory.lock.Unlock()

	redisHandler, exists = factory.redisMap[channel]
	if !exists {
		redisResolver, exists := factory.redisResolverMap[channel]
		if !exists {
			return nil, errors.New("redis channel " + channel + " not exists")
		}

		var err error = nil
		redisHandler, err = redisResolver()
		if err != nil {
			return nil, errors.New("redis resolve fail, channel:" + channel + ", error:" + err.Error())
		}
		factory.redisMap[channel] = redisHandler
	}

	return redisHandler, nil
}

func (factory *Factory) MakeRedis(config Config) (redis.Cmdable, error) {
	fields := helper.ValidateAndGetErrFields(config)
	if len(fields) > 0 {
		return nil, errors.New("redis config error, reason: fields: " + strings.Join(fields, ","))
	}

	return redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(int(config.Port)),
		Username: config.Username,
		Password: config.Password,
		DB:       int(config.Db),
		PoolSize: int(config.PoolSize),
	}), nil
}

func (factory *Factory) RegisterRedis(channel string, redisResolver func() (redis.Cmdable, error)) {
	factory.redisResolverMap[channel] = redisResolver
}

func (factory *Factory) Register(maps map[string]Config) {
	for key, value := range maps {
		func(channel string, config Config) {
			factory.RegisterRedis(channel, func() (redis.Cmdable, error) {
				return factory.MakeRedis(config)
			})
		}(key, value)
	}
}
