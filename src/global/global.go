package global

import (
	"github.com/asaskevich/EventBus"
	ut "github.com/go-playground/universal-translator"
	"github.com/we7coreteam/w7-rangine-go/src/components/database"
	"github.com/we7coreteam/w7-rangine-go/src/components/redis"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger"
)

type RGlobal struct {
	LoggerFactory   *logger.LoggerFactory
	Event           EventBus.Bus
	RedisFactory    *redis.RedisFactory
	DatabaseFactory *database.DatabaseFactory
	Translator      ut.Translator
}

var G = new(RGlobal)
